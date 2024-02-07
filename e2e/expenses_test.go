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

type basicExpenseInfo struct {
	name          string
	description   string
	paymentMethod string
	amount        int
	currency      string
}

func getDummyExpenseInfo() basicExpenseInfo {
	return basicExpenseInfo{
		name:          "name",
		description:   "description",
		paymentMethod: "card",
		amount:        1500,
		currency:      "RON",
	}
}

func compareExpenseWithResponse(t *testing.T, expenseInfo basicExpenseInfo, responseData casheerapi.ExpenseData) {
	if responseData.Type != "expense" {
		t.Errorf("Invalid return type, expected \"expense\" but got \"%s\"\n", responseData.Type)
	}

	// Check attributes
	if responseData.Attributes.Name != expenseInfo.name {
		t.Errorf("Expense name doesn't match; expected %s but got %s\n", expenseInfo.name, responseData.Attributes.Name)
	}
	if responseData.Attributes.Description != expenseInfo.description {
		t.Errorf("Expense description doesn't match; expected %s but got %s\n", expenseInfo.description, responseData.Attributes.Description)
	}
	if responseData.Attributes.PaymentMethod != expenseInfo.paymentMethod {
		t.Errorf("Expense payment method doesn't match; expected %s but got %s\n", expenseInfo.paymentMethod, responseData.Attributes.PaymentMethod)
	}

	val := currency.NewRONValue(expenseInfo.amount)
	if responseData.Attributes.Value != (casheerapi.MonetaryValueAttributes{
		Amount:   val.Amount,
		Exponent: val.Exponent,
		Currency: val.Currency,
	}) {
		t.Errorf("Expense value doesn't match; expected %v but got %v\n", responseData.Attributes.Value, casheerapi.MonetaryValueAttributes{
			Amount:   val.Amount,
			Exponent: val.Exponent,
			Currency: val.Currency,
		})
	}
}

func checkExpenseLinks(t *testing.T, entid int, responseData casheerapi.ExpenseData) {
	if !strings.Contains(responseData.Links.Self, fmt.Sprintf("%d/expenses/%s", entid, responseData.Id)) {
		t.Errorf("Invalid related link format; got %s\n", responseData.Links.Self)
	}

	// Check related resources
	if !strings.Contains(responseData.Relationships.Entries.Links.Related, strconv.Itoa(entid)) {
		t.Errorf("Expected related link (%s) to include expense id (%d)\n", responseData.Relationships.Entries.Links.Related, entid)
	}
}

func Test_CreateBasicExpense_ExpenseIsCreated_and_ReturnedValuesAreCorrect(t *testing.T) {

	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})
	entid := setupEntry(t)

	expenseInfo := getDummyExpenseInfo()
	expenseResponse, err := casheerClient.CreateBasicExpense(entid, expenseInfo.name, expenseInfo.description, expenseInfo.paymentMethod, expenseInfo.amount, expenseInfo.currency)
	if err != nil {
		t.Fatalf("Did not expect error when creating expense, but got error: %s\n", err.Error())
	}

	compareExpenseWithResponse(t, expenseInfo, expenseResponse.Data)
	checkExpenseLinks(t, entid, expenseResponse.Data)
}

func Test_DeleteExpense_ExistingExpenseIsDeleted(t *testing.T) {

	t.Cleanup(func() {
		store.DeleteAllData(conn)
	})
	entid := setupEntry(t)

	expenseInfo := getDummyExpenseInfo()
	expenseResponseCreate, err := casheerClient.CreateBasicExpense(entid, expenseInfo.name, expenseInfo.description, expenseInfo.paymentMethod, expenseInfo.amount, expenseInfo.currency)
	if err != nil {
		t.Fatalf("Did not expect error when creating expense, but got error: %s\n", err.Error())
	}

	expid, _ := strconv.Atoi(expenseResponseCreate.Data.Id)
	expenseResponseDelete, err := casheerClient.DeleteExpenseForEntry(entid, expid)
	if err != nil {
		t.Fatalf("Did not expect error when deleting existing expense, but got error: %s\n", err.Error())
	}
	// Do not check links because a deleted expense should not have links.
	compareExpenseWithResponse(t, expenseInfo, expenseResponseDelete.Data)

}
