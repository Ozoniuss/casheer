package e2e

import (
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/internal/store"
)

func TestCreateEntryWithExpensesFlow(t *testing.T) {

	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})

	createResp, err := casheerClient.CreateEntry(10, 2023, "category", "subcategory", 1000, "EUR", false)
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

	_, err = casheerClient.CreateEntry(10, 2023, "category", "subcategory", 1000, "EUR", false)
	if err == nil {
		t.Fatal("Idempotency criteria violated; entry with the same identifying information created again.")
	}

	expenseResp, err := casheerClient.CreateBasicExpense(entid, "car trip", "big fuel", "card", 1000, "RON")
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
		expenseResp.Data.Attributes.Value.Amount != 1000 ||
		expenseResp.Data.Attributes.Value.Currency != "RON" ||
		expenseResp.Data.Attributes.Value.Exponent != -2 {
		t.Errorf("Received invalid debt attributes after creating the debt.")
	}
	_, err = casheerClient.CreateBasicExpenseWithoutId("category", "subcategory", 10, 2023, "second car trip", "more fuel", "card", 500, "RON")
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

func TestCreateListFilterEntryFlow(t *testing.T) {

	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})

	createRespEntry1, err := casheerClient.CreateEntry(10, 2022, "category1", "subcategory1", 5000, "EUR", true)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created entry with id %s\n", createRespEntry1.Data.Id)
	createRespEntry2, err := casheerClient.CreateEntry(10, 2022, "category1", "subcategory2", 5000, "EUR", true)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created entry with id %s\n", createRespEntry2.Data.Id)
	createRespEntry3, err := casheerClient.CreateEntry(10, 2022, "category2", "subcategory2", 5000, "EUR", true)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created entry with id %s\n", createRespEntry3.Data.Id)

	allEntries, err := casheerClient.ListEntries()

	if err != nil {
		t.Error("Expected no error when listing entries.")
	}

	if total := len(allEntries.Data); total != 3 {
		t.Errorf("Expected to have 3 entries, got %d\n", total)
	}
}
