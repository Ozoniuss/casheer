package e2e

import (
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

var casheerEntriesClient, _ = httpclient.NewCasheerHTTPClient(
	httpclient.WithCustomAuthority("localhost:6597"),
)

func TestCreateEntryWithExpensesFlow(t *testing.T) {
	createResp, err := casheerEntriesClient.CreateEntry(10, 2023, "category", "subcategory", 1000, "EUR", false)
	if err != nil {
		t.Fatalf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created entry with id %s\n", createResp.Data.Id)

	entid, err := strconv.Atoi(createResp.Data.Id)
	if err != nil {
		t.Fatalf("Could not convert entry id to int: %s\n", err.Error())
	}

	if createResp.Data.Attributes.Category != "category" ||
		createResp.Data.Attributes.Subcategory != "subcategory" ||
		createResp.Data.Attributes.ExpectedTotal.Amount != 1000 ||
		createResp.Data.Attributes.ExpectedTotal.Currency != "EUR" ||
		createResp.Data.Attributes.ExpectedTotal.Exponent != -2 ||
		createResp.Data.Attributes.Month != 10 ||
		createResp.Data.Attributes.Year != 2023 ||
		createResp.Data.Attributes.Recurring != false {
		t.Errorf("Received invalid debt attributes after creating the debt.")
	}

	_, err = casheerEntriesClient.CreateEntry(10, 2023, "category", "subcategory", 1000, "EUR", false)
	if err == nil {
		t.Fatal("Idempotency criteria violated; entry with the same identifying information created again.")
	}

	expenseResp, err := casheerEntriesClient.CreateBasicExpense(entid, "car trip", "big fuel", "card", 1000, "RON")
	if err != nil {
		t.Fatalf("Did not expect error when creating expense; got %s\n", err.Error())
	}

	_, err = strconv.Atoi(expenseResp.Data.Id)
	if err != nil {
		t.Fatalf("Could not convert expense id to int: %s\n", err.Error())
	}

	if expenseResp.Data.Attributes.Name != "car trip" ||
		expenseResp.Data.Attributes.Description != "big fuel" ||
		expenseResp.Data.Attributes.PaymentMethod != "card" ||
		expenseResp.Data.Attributes.Amount != 1000 ||
		expenseResp.Data.Attributes.Currency != "RON" ||
		expenseResp.Data.Attributes.Exponent != -2 {
		t.Errorf("Received invalid debt attributes after creating the debt.")
	}
	_, err = casheerEntriesClient.CreateBasicExpenseWithoutId("category", "subcategory", 10, 2023, "second car trip", "more fuel", "card", 500, "RON")
	if err != nil {
		t.Fatalf("Did not expect error when creating expense; got %s\n", err.Error())
	}

	// getResp, err := casheerDebtClient.GetDebt(did)
	// if err != nil {
	// 	t.Error("Expected the returned debt Id to point to a valid debt, got error instead.")
	// }

	// if createResp.Data != getResp.Data {
	// 	t.Errorf("Debt returned in GET request should match debt returned from POST request. Got %+v, want %+v", createResp.Data, getResp.Data)
	// }

	// deleteResp, err := casheerDebtClient.DeleteDebt(did)
	// if err != nil {
	// 	t.Error("Expected to be able to delete the newly created debt.")
	// }

	// if deleteResp.Data.Attributes != getResp.Data.Attributes {
	// 	t.Errorf("Attributes of debt returned in DELETE request should match debt returned from POST request. Got %+v, want %+v", createResp.Data, getResp.Data)
	// }
	// if deleteResp.Data.Links.Self != "" {
	// 	t.Errorf("Expected no self link, got %s\n", deleteResp.Data.Links.Self)
	// }
	// t.Logf("Deleted debt with id %d\n", did)

	// _, err = casheerDebtClient.GetDebt(did)
	// if err == nil {
	// 	t.Fatal("Expected to get an error when retrieving deleted debt, got none.")
	// }

	// jsonerr, ok := err.(casheerapi.ErrorResponse)
	// if !ok {
	// 	t.Fatalf("Expected to get error of type jsonapi, got %v\n", err)
	// }

	// if jsonerr.Err.Status != 404 {
	// 	t.Errorf("Expected to get status 404, got %d\n", jsonerr.Err.Status)
	// }
}

// func TestCreateListFilterFlow(t *testing.T) {
// 	createRespMarian1, err := casheerDebtClient.CreateDebt("Marian", "get tf out", 100, "EUR", -2)
// 	if err != nil {
// 		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
// 	}
// 	t.Logf("Created debt with id %s\n", createRespMarian1.Data.Id)

// 	createRespMarian2, err := casheerDebtClient.CreateDebt("Marian", "get tf out", 200, "EUR", -2)
// 	if err != nil {
// 		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
// 	}
// 	t.Logf("Created debt with id %s\n", createRespMarian2.Data.Id)

// 	createRespDaniel, err := casheerDebtClient.CreateDebt("Daniels", "get tf out", 50, "EUR", -2)
// 	if err != nil {
// 		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
// 	}
// 	t.Logf("Created debt with id %s\n", createRespDaniel.Data.Id)

// 	allDebts, err := casheerDebtClient.ListDebts()

// 	if err != nil {
// 		t.Error("Expected no error when listing debts.")
// 	}

// 	if total := len(allDebts.Data); total != 3 {
// 		t.Errorf("Expected to have 3 debts, got %d\n", total)
// 	}

// 	MariansDebts, err := casheerDebtClient.ListDebtsForPerson("Marian")

// 	if err != nil {
// 		t.Error("Expected no error when listing debts.")
// 	}

// 	if total := len(MariansDebts.Data); total != 2 {
// 		t.Errorf("Expected to have 2 debts for Marian, got %d\n", total)
// 	}
// }
