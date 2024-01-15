package e2e

import (
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/internal/store"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func TestCreateRetrieveDeleteDebtFlow(t *testing.T) {
	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})
	createResp, err := casheerClient.CreateDebt("Marian", "get tf out", 100, "EUR", -2)
	if err != nil {
		t.Fatalf("Did not expect error, but got error: %s\n", err.Error())
	}
	if createResp.Data.Attributes.Person != "Marian" {
		t.Errorf("Debt does not match with what was created, got response: %+v\n", createResp)
	}
	t.Logf("Created debt with id %s\n", createResp.Data.Id)

	did, err := strconv.Atoi(createResp.Data.Id)
	if err != nil {
		t.Fatalf("Could not convert debt id to int: %s\n", err.Error())
	}

	if createResp.Data.Attributes.Person != "Marian" ||
		createResp.Data.Attributes.Value.Currency != "EUR" ||
		createResp.Data.Attributes.Value.Amount != 100 ||
		createResp.Data.Attributes.Value.Exponent != -2 ||
		createResp.Data.Attributes.Details != "get tf out" {
		t.Errorf("Received invalid debt attributes after creating the debt.")
	}

	getResp, err := casheerClient.GetDebt(did)
	if err != nil {
		t.Error("Expected the returned debt Id to point to a valid debt, got error instead.")
	}

	if createResp.Data != getResp.Data {
		t.Errorf("Debt returned in GET request should match debt returned from POST request. Got %+v, want %+v", createResp.Data, getResp.Data)
	}

	deleteResp, err := casheerClient.DeleteDebt(did)
	if err != nil {
		t.Error("Expected to be able to delete the newly created debt.")
	}

	if deleteResp.Data.Attributes != getResp.Data.Attributes {
		t.Errorf("Attributes of debt returned in DELETE request should match debt returned from POST request. Got %+v, want %+v", createResp.Data, getResp.Data)
	}
	if deleteResp.Data.Links.Self != "" {
		t.Errorf("Expected no self link, got %s\n", deleteResp.Data.Links.Self)
	}
	t.Logf("Deleted debt with id %d\n", did)

	_, err = casheerClient.GetDebt(did)
	if err == nil {
		t.Fatal("Expected to get an error when retrieving deleted debt, got none.")
	}

	jsonerr, ok := err.(casheerapi.ErrorResponse)
	if !ok {
		t.Fatalf("Expected to get error of type jsonapi, got %v\n", err)
	}

	if jsonerr.Err.Status != 404 {
		t.Errorf("Expected to get status 404, got %d\n", jsonerr.Err.Status)
	}
}

func TestCreateListFilterDebtFlow(t *testing.T) {

	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})

	createRespMarian1, err := casheerClient.CreateDebt("Marian", "get tf out", 100, "EUR", -2)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created debt with id %s\n", createRespMarian1.Data.Id)

	createRespMarian2, err := casheerClient.CreateDebt("Marian", "get tf out", 200, "EUR", -2)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created debt with id %s\n", createRespMarian2.Data.Id)

	createRespDaniel, err := casheerClient.CreateDebt("Daniels", "get tf out", 50, "EUR", -2)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
	}
	t.Logf("Created debt with id %s\n", createRespDaniel.Data.Id)

	allDebts, err := casheerClient.ListDebts()

	if err != nil {
		t.Error("Expected no error when listing debts.")
	}

	if total := len(allDebts.Data); total != 3 {
		t.Errorf("Expected to have 3 debts, got %d\n", total)
	}

	MariansDebts, err := casheerClient.ListDebtsForPerson("Marian")

	if err != nil {
		t.Error("Expected no error when listing debts.")
	}

	if total := len(MariansDebts.Data); total != 2 {
		t.Errorf("Expected to have 2 debts for Marian, got %d\n", total)
	}
}
