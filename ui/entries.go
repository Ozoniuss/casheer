package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Ozoniuss/casheer/client/httpclient"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func createEntryWithCasheerClient(cl *httpclient.CasheerHTTPClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

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
}
