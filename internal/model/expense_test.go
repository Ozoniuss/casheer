package model

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidExpense(t *testing.T) {
	tests := []struct {
		id      int
		expense Expense
		valid   bool
	}{
		{
			id: 1,
			expense: Expense{
				EntryId: 1,
				Value: Value{
					Currency: "abcd",
					Amount:   1000,
					Exponent: -2,
				},
				Name: "pizza",
			},
			valid: false,
		},
		{
			id: 2,
			expense: Expense{
				EntryId: 2,
				Value: Value{
					Currency: "EUR",
					Amount:   1000,
					Exponent: -2,
				},
				Name: "pizza",
			},
			valid: true,
		},
	}
	validate := validator.New()
	for _, tt := range tests {
		err := validate.Struct(tt.expense)
		if tt.valid {
			if err != nil {
				t.Fatalf("test %d failed: got unwanted validation error for expense: %s", tt.id, err.Error())
			}
		} else {
			if err == nil {
				t.Fatalf("test %d failed: expense should not be valid", tt.id)
			}
		}
	}
}

func TestValidExpenseNewValidation(t *testing.T) {
	tests := []struct {
		name    string
		expense Expense
		valid   bool
	}{
		{
			name: "valid expense",
			expense: Expense{
				Name: "pizza",
			},
			valid: true,
		},
		{
			name: "empty expense name",
			expense: Expense{
				Name: "",
			},
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.expense.Validate() == nil) != tt.valid {
				t.Fatalf("expected to be valid: %t, got validation error: %v", tt.valid, tt.expense.Validate())
			}
		})
	}
}
