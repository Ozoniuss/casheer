package main

import (
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"text/template"

	"github.com/Ozoniuss/casheer/client/httpclient"
	"golang.org/x/exp/maps"
)

type TemplateData struct {
	DebtsList            []DebtListItem
	CategorizedEntryList []CategoryWithEntries
}

func main() {

	cl, err := httpclient.NewCasheerHTTPClient()
	if err != nil {
		panic(err)
	}

	// cl.CreateDebt("Marian", "Ce pula mea", 100, "RON", -2)

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		debts := loadDebtsList(cl)
		entries := loadCategorizedEntriesList(cl)
		data := TemplateData{
			DebtsList:            debts,
			CategorizedEntryList: entries,
		}
		tmpl.Execute(w, data)
	}

	handleDeleteDebt := func(w http.ResponseWriter, r *http.Request) {

		didstr := r.PostFormValue("debtid")
		did, err := strconv.Atoi(didstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid debt id: %s", err.Error())
			return
		}

		_, err = cl.DeleteDebt(did)
		if err != nil {
			w.WriteHeader(r.Response.StatusCode)
			fmt.Fprintf(w, "could not delete debt: %s", err.Error())
		}
		w.WriteHeader(http.StatusOK)
	}

	handleCreateDebt := func(w http.ResponseWriter, r *http.Request) {

		person := r.FormValue("person")
		totalMoneyStr := r.FormValue("total-money")
		currency := r.FormValue("currency")
		details := r.FormValue("details")

		if person == "" || currency == "" {
			w.WriteHeader(400)
			fmt.Fprint(w, "some fields should not be empty")
			return
		}

		totalMoney, err := strconv.Atoi(totalMoneyStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid money value: %s", err.Error())
			return
		}

		resp, err := cl.CreateDebt(person, details, totalMoney, currency, -2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not create debt: %s", err.Error())
			return
		}

		did, _ := strconv.Atoi(resp.Data.Id)

		dli := DebtListItem{
			Id:         did,
			Person:     resp.Data.Attributes.Person,
			TotalMoney: float32(math.Pow10(resp.Data.Attributes.Value.Exponent)) * float32(resp.Data.Attributes.Value.Amount),
			Currency:   resp.Data.Attributes.Value.Currency,
			Details:    resp.Data.Attributes.Details,
		}

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "debt-list-element", dli)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/deleteDebt", handleDeleteDebt)
	http.HandleFunc("/createDebt", handleCreateDebt)

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
	entries, err := c.ListEntries()
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
		}

		for _, exp := range e2.Expenses {
			if exp.Currency == "RON" {
				e2.RunningTotal += exp.TotalMoney
			}
		}

		data = append(data, e2)
	}
	return createCategoriesArray(data)
}

func createCategoriesArray(entries []EntryListItem) []CategoryWithEntries {
	categories := make(map[string][]EntryListItem)
	for _, e := range entries {
		categories[e.Category] = append(categories[e.Category], e)
		// if _, ok := categories[e.Category]; ok {
		// 	categories[e.Category] = append(categories[e.Category], e)
		// } else {
		// 	categories[e.Category] = []EntryListItem{e}
		// }
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
	entries, err := c.GetEntry(entryId)
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
