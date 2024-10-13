package main

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"text/template"
	"time"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

func createExpenseWithCasheerClient(cl *httpclient.CasheerHTTPClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

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
}

func deleteExpenseWithCasheerClient(cl *httpclient.CasheerHTTPClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
