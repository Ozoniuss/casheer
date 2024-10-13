package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"slices"
	"strconv"
	"text/template"
	"time"

	"github.com/Ozoniuss/casheer/client/httpclient"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
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

// debugging
// var year = 2023
// var month = 12

func main() {

	cl, err := httpclient.NewCasheerHTTPClient()
	if err != nil {
		panic(err)
	}

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		if err != nil {
			panic(err)
		}
		debts := loadDebtsList(cl)
		entries := loadCategorizedEntriesList(cl)
		templateData = TemplateData{
			DebtsList:            debts,
			CategorizedEntryList: entries,
		}
		err = tmpl.Execute(w, templateData)
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
	handleDeleteExpense := func(w http.ResponseWriter, r *http.Request) {
		entidstr := r.PostFormValue("entid")
		entid, err := strconv.Atoi(entidstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid entry id: %s", err.Error())
			return
		}
		expidstr := r.PostFormValue("expid")
		expid, err := strconv.Atoi(expidstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid expense id: %s", err.Error())
			return
		}

		_, err = cl.DeleteExpenseForEntry(entid, expid)
		if err != nil {
			w.WriteHeader(r.Response.StatusCode)
			fmt.Fprintf(w, "could not delete expense: %s", err.Error())
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

		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "debt-list-element", dli)
	}

	handleCreateExpense := func(w http.ResponseWriter, r *http.Request) {

		name := r.FormValue("name")
		category := r.FormValue("category")
		subcategory := r.FormValue("subcategory")
		totalMoneyStr := r.FormValue("total-money")
		currency := r.FormValue("currency")
		description := r.FormValue("details")
		paymentMethod := r.FormValue("payment-method")

		// basic validation
		if subcategory == "" || currency == "" {
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

		// category is given automatically, search the id in the list
		// of expenses
		foundCategoryIndex := slices.IndexFunc[[]CategoryWithEntries, CategoryWithEntries](templateData.CategorizedEntryList, func(cwe CategoryWithEntries) bool {
			return cwe.Category == category
		})
		if foundCategoryIndex == -1 {
			panic("did not return category correctly. " + category)
		}

		catWithEntries := templateData.CategorizedEntryList[foundCategoryIndex]
		entryIdx := slices.IndexFunc[[]EntryListItem, EntryListItem](catWithEntries.Entries, func(eli EntryListItem) bool {
			return eli.Subcategory == subcategory
		})
		// In the categorized list, we find the entry to which this shit
		// belongs.
		var entryId int
		if entryIdx != -1 {
			// we have an existing entry for this
			entryId = catWithEntries.Entries[entryIdx].Id
		} else {
			// we need a new entry
			entryId = -1
		}

		// var createdEntry *casheerapi.CreateEntryResponse
		if entryId == -1 {
			resp, err := cl.CreateEntry(
				int(time.Now().Month()),
				time.Now().Year(),
				category,
				subcategory,
				totalMoney,
				currency,
				false,
			)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "could not create entry: %s", err.Error())
				return
			}
			entryId, err = strconv.Atoi(resp.Data.Id)
			if err != nil {
				panic("something went wrong")
			}

			// keep track if entry was created
			// createdEntry = &resp
		}

		_, err = cl.CreateBasicExpense(
			entryId,
			name,
			description,
			paymentMethod,
			totalMoney,
			currency,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not create expense: %s", err.Error())
			return
		}

		templateData.CategorizedEntryList = loadCategorizedEntriesList(cl)

		idx := slices.IndexFunc[[]CategoryWithEntries, CategoryWithEntries](templateData.CategorizedEntryList, func(cwe CategoryWithEntries) bool {
			return category == cwe.Category
		})

		fmt.Println(templateData.CategorizedEntryList[idx])

		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "all-entries-categorized", templateData.CategorizedEntryList[idx])
	}

	handleCreateEntry := func(w http.ResponseWriter, r *http.Request) {

		category := r.FormValue("category")
		subcategory := r.FormValue("subcategory")
		expectedTotalStr := r.FormValue("expected-total")
		recurringCheck := r.FormValue("recurring")
		currency := r.FormValue("currency")
		monthStr := r.FormValue("month")
		yearStr := r.FormValue("year")

		recurring := false
		if recurringCheck == "on" {
			recurring = true
		}

		// basic validation
		if subcategory == "" || category == "" {
			w.WriteHeader(400)
			fmt.Fprint(w, "some fields should not be empty")
			return
		}

		if currency == "" {
			currency = "RON"
		}

		expectedTotal, err := strconv.Atoi(expectedTotalStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid expected total: %s", err.Error())
			return
		}
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid month: %s", err.Error())
			return
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid year: %s", err.Error())
			return
		}

		createdEntry, err := cl.CreateEntry(month, year, category, subcategory, expectedTotal, currency, recurring)
		if err != nil {
			var apierr casheerapi.ErrorResponse
			if errors.As(err, &apierr) {
				fmt.Printf("got error %v\n", apierr)
				w.WriteHeader(apierr.Err.Status)
				fmt.Fprintf(w, "could not create entry: %s", apierr.Err.Detail)
			} else {
				fmt.Printf("got error %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "could not create entry: %s", err)
			}
		}

		fmt.Printf("created entry %+v\n", createdEntry)

		templateData.CategorizedEntryList = loadCategorizedEntriesList(cl)

		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "all-entries-categorized-inner", templateData)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/deleteDebt", handleDeleteDebt)
	http.HandleFunc("/deleteExpense", handleDeleteExpense)
	http.HandleFunc("/createDebt", handleCreateDebt)
	http.HandleFunc("/createExpense", handleCreateExpense)
	http.HandleFunc("/createEntry", handleCreateEntry)
	http.HandleFunc("/year", func(w http.ResponseWriter, r *http.Request) {
		// year=2023
		yearstr, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		year, _ = strconv.Atoi(string(yearstr)[5:])

		// This is so ugly I almost want to switch to being a java developer.
		// A bigass cleanup is coming up soon when I'll think about how to
		// organize this mess.
		debts := loadDebtsList(cl)
		entries := loadCategorizedEntriesList(cl)
		templateData = TemplateData{
			DebtsList:            debts,
			CategorizedEntryList: entries,
		}
		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "all-planning-data", templateData)

		fmt.Printf("changing year to %d\n", year)
	})
	http.HandleFunc("/month", func(w http.ResponseWriter, r *http.Request) {
		// month=11
		monthstr, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		month, _ = strconv.Atoi(string(monthstr)[6:])

		debts := loadDebtsList(cl)
		entries := loadCategorizedEntriesList(cl)
		templateData = TemplateData{
			DebtsList:            debts,
			CategorizedEntryList: entries,
		}
		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "all-planning-data", templateData)

		fmt.Printf("changing month to %d\n", month)
	})
	// TODO: use new http library for setting the GET method
	http.HandleFunc("/period", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]int{
			"year":  year,
			"month": month,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

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
