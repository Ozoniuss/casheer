package main

import (
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"text/template"
	"time"

	"github.com/Ozoniuss/casheer/client/httpclient"
	"golang.org/x/exp/maps"
)

type TemplateData struct {
	DebtsList            []DebtListItem
	CategorizedEntryList []CategoryWithEntries
}

var Funcs template.FuncMap = template.FuncMap{
	"Sum": func(x int, y float32) int {
		return int(math.Round(float64(x) + float64(y)))
	},
}

var templateData TemplateData
var allEntries []EntryListItem

var year = time.Now().Year()
var month = int(time.Now().Month())

func main() {

	cl, err := httpclient.NewCasheerHTTPClient()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", indexHandlerWithCasheerClient(cl))
	http.HandleFunc("/deleteDebt", deleteDebtWithCasheerClient(cl))
	http.HandleFunc("/deleteExpense", deleteExpenseWithCasheerClient(cl))
	http.HandleFunc("/createDebt", createDebtWithCasheerClient(cl))
	http.HandleFunc("/createExpense", createExpenseWithCasheerClient(cl))
	http.HandleFunc("/createEntry", createEntryWithCasheerClient(cl))
	http.HandleFunc("/year", yearHandlerWithCasheerClient(cl))
	http.HandleFunc("/month", monthHandlerWithCasheerClient(cl))
	http.HandleFunc("/period", periodHandler)

	http.ListenAndServe(":7145", nil)
}

func loadDebtsList(c *httpclient.CasheerHTTPClient) []DebtListItem {
	debts, err := c.ListDebts()
	if err != nil {
		panic(err)
	}
	data := []DebtListItem{}
	for _, d := range debts.Data {
		did, _ := strconv.Atoi(d.Id)
		d2 := DebtListItem{
			Id:         did,
			Person:     d.Attributes.Person,
			TotalMoney: float32(math.Pow10(d.Attributes.Value.Exponent)) * float32(d.Attributes.Value.Amount),
			Currency:   d.Attributes.Value.Currency,
			Details:    d.Attributes.Details,
		}
		data = append(data, d2)
	}
	return data
}

func loadCategorizedEntriesList(c *httpclient.CasheerHTTPClient) []CategoryWithEntries {
	entries, err := c.ListEntriesForPeriod(month, year)
	if err != nil {
		panic(err)
	}
	data := []EntryListItem{}
	for _, e := range entries.Data {
		eid, _ := strconv.Atoi(e.Id)
		e2 := EntryListItem{
			Id:          eid,
			TotalMoney:  float32(math.Pow10(e.Attributes.ExpectedTotal.Exponent)) * float32(e.Attributes.ExpectedTotal.Amount),
			Currency:    e.Attributes.ExpectedTotal.Currency,
			Category:    e.Attributes.Category,
			Subcategory: e.Attributes.Subcategory,
			Recurring:   e.Attributes.Recurring,
			Expenses:    loadExpensesListForEntry(c, eid),
			// hardcoding for now since I don't know htmx any better
			RunningTotal: map[string]float32{"EUR": 0, "RON": 0, "USD": 0, "GBP": 0},
		}

		for _, exp := range e2.Expenses {
			if exp.Currency != "RON" && exp.Currency != "EUR" && exp.Currency != "USD" && exp.Currency != "GBP" {
				panic("unsupported currency " + exp.Currency)
			}
			e2.RunningTotal[exp.Currency] += exp.TotalMoney
		}

		data = append(data, e2)
	}
	// store this globally so it's accessible in other places as well
	allEntries = data
	fmt.Println(allEntries)
	return createCategoriesArray(data)
}

func createCategoriesArray(entries []EntryListItem) []CategoryWithEntries {
	categories := make(map[string][]EntryListItem)
	for _, e := range entries {
		categories[e.Category] = append(categories[e.Category], e)
	}

	categoriesWithEntries := make([]CategoryWithEntries, 0)

	// place income first
	if _, ok := categories["income"]; ok {
		categoriesWithEntries = append(categoriesWithEntries, CategoryWithEntries{
			Category: "income",
			Entries:  categories["income"],
		})
	}

	// place all other after that, sorted
	categoriesSorted := maps.Keys(categories)
	slices.Sort(categoriesSorted)

	for _, c := range categoriesSorted {
		if c == "income" {
			continue
		}
		categoriesWithEntries = append(categoriesWithEntries, CategoryWithEntries{
			Category: c,
			Entries:  categories[c],
		})
	}
	return categoriesWithEntries
}

func loadExpensesListForEntry(c *httpclient.CasheerHTTPClient, entryId int) []ExpenseListItem {
	entries, err := c.GetEntry(entryId, true)
	if err != nil {
		panic(err)
	}
	data := []ExpenseListItem{}

	if entries.Included == nil {
		panic("wtf")
	}

	for _, exp := range *entries.Included {
		if exp.Type != "expense" {
			fmt.Println("debug")
			continue
		}

		eid, _ := strconv.Atoi(exp.Id)
		d2 := ExpenseListItem{
			Id:            eid,
			Name:          exp.Attributes.Name,
			TotalMoney:    float32(math.Pow10(exp.Attributes.Value.Exponent)) * float32(exp.Attributes.Value.Amount),
			Currency:      exp.Attributes.Value.Currency,
			PaymentMethod: exp.Attributes.PaymentMethod,
			Description:   exp.Attributes.Description,
		}
		data = append(data, d2)
	}
	return data
}
