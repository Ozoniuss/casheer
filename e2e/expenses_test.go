package e2e

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Ozoniuss/casheer/currency"
	"github.com/Ozoniuss/casheer/internal/store"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func setupEntry(t *testing.T) int {
	id, err := store.CreateDummyEntry(conn)
	if err != nil {
		t.Fatalf("failed creating dummy entry: %s", err.Error())
	}
	t.Logf("created dummy entry with id %d", id)
	return id
}

func Test_CreateBasicExpense_ExpenseIsCreated_and_ReturnedValuesAreCorrect(t *testing.T) {

	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})
	entid := setupEntry(t)

	var name = "expense 1"
	var description = "description"
	var paymentMethod = "card"
	var amount = 1500
	var ccurrency = "RON"
	expenseResponse, err := casheerClient.CreateBasicExpense(entid, name, description, paymentMethod, amount, ccurrency)
	if err != nil {
		t.Fatalf("Did not expect error when creating expense, but got error: %s\n", err.Error())
	}

	if expenseResponse.Data.Type != "expense" {
		t.Errorf("Invalid return type, expected \"expense\" but got \"%s\"\n", expenseResponse.Data.Type)
	}

	// Check attributes
	if expenseResponse.Data.Attributes.Name != name {
		t.Errorf("Expense name doesn't match; expected %s but got %s\n", name, expenseResponse.Data.Attributes.Name)
	}
	if expenseResponse.Data.Attributes.Description != description {
		t.Errorf("Expense description doesn't match; expected %s but got %s\n", description, expenseResponse.Data.Attributes.Description)
	}
	if expenseResponse.Data.Attributes.PaymentMethod != paymentMethod {
		t.Errorf("Expense payment method doesn't match; expected %s but got %s\n", paymentMethod, expenseResponse.Data.Attributes.PaymentMethod)
	}
	val := currency.NewRONValue(1500)
	if expenseResponse.Data.Attributes.Value != (casheerapi.MonetaryValueAttributes{
		Amount:   val.Amount,
		Exponent: val.Exponent,
		Currency: val.Currency,
	}) {
		t.Errorf("Expense value doesn't match; expected %v but got %v\n", expenseResponse.Data.Attributes.Value, casheerapi.MonetaryValueAttributes{
			Amount:   val.Amount,
			Exponent: val.Exponent,
			Currency: val.Currency,
		})
	}

	if !strings.Contains(expenseResponse.Data.Links.Self, fmt.Sprintf("%d/expenses/%s", entid, expenseResponse.Data.Id)) {
		t.Errorf("Invalid related link format; got %s\n", expenseResponse.Data.Links.Self)
	}

	// Check related resources
	if !strings.Contains(expenseResponse.Data.Relationships.Entries.Links.Related, strconv.Itoa(entid)) {
		t.Errorf("Expected related link (%s) to include expense id (%d)\n", expenseResponse.Data.Relationships.Entries.Links.Related, entid)
	}
}
