package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

func indexHandlerWithCasheerClient(cl *httpclient.CasheerHTTPClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index.html").Funcs(Funcs).ParseFiles("index.html"))
		debts := loadDebtsList(cl)
		entries := loadCategorizedEntriesList(cl)
		templateData = TemplateData{
			DebtsList:            debts,
			CategorizedEntryList: entries,
		}
		err := tmpl.Execute(w, templateData)
		if err != nil {
			// todo: render an error template instead
			panic(err)
		}
	}
}

func yearHandlerWithCasheerClient(cl *httpclient.CasheerHTTPClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
func monthHandlerWithCasheerClient(cl *httpclient.CasheerHTTPClient) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

var periodHandler = func(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]int{
		"year":  year,
		"month": month,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
