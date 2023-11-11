package model

import (
	"testing"
)

func TestVadidDebt(t *testing.T) {
	tests := []struct {
		name  string
		debt  Debt
		valid bool
	}{
		{
			name: "valid debt",
			debt: Debt{
				Person:  "Andrei",
				Details: "pay me back",
			},
			valid: true,
		},
		{
			name: "empty person",
			debt: Debt{
				Person:  "",
				Details: "pay me back",
			},
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.debt.Validate() == nil) != tt.valid {
				t.Fatalf("expected to be valid: %t, got validation error: %v", tt.valid, tt.debt.Validate())
			}
		})
	}
}
