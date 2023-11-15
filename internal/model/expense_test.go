package model

import (
	"testing"
)

func TestValidExpense(t *testing.T) {
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
