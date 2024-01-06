package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

func main() {

	cl, err := httpclient.NewCasheerHTTPClient()
	if err != nil {
		panic(err)
	}

	// cl.CreateDebt("Marian", "Ce pula mea", 100, "RON", -2)

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		data := loadDebtsList(cl)
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
