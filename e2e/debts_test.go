package e2e

import (
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/client/httpclient"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

var casheerDebtClient, _ = httpclient.NewCasheerHTTPClient()

func TestCreateRetrieveDeleteDebtFlow(t *testing.T) {
	createResp, err := casheerDebtClient.CreateDebt("Marian", "get tf out", 100, "EUR", -2)
	if err != nil {
		t.Errorf("Did not expect error, but got error: %s\n", err.Error())
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
		createResp.Data.Attributes.Currency != "EUR" ||
		createResp.Data.Attributes.Amount != 100 ||
		createResp.Data.Attributes.Exponent != -2 ||
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
