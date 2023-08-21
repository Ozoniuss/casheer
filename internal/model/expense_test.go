package model

import (
	"testing"

	"github.com/Ozoniuss/casheer/internal/currency"
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
				Value: currency.Value{
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
				Value: currency.Value{
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
