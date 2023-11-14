package model

import (
	"testing"
)

func TestValidEntryNewValidation(t *testing.T) {
	tests := []struct {
		name  string
		entry Entry
		valid bool
	}{
		{
			name: "valid entry",
			entry: Entry{
				Month:       10,
				Year:        2022,
				Category:    "one",
				Subcategory: "two",
			},
			valid: true,
		},
		{
			name: "invalid month",
			entry: Entry{
				Month:       13,
				Year:        2022,
				Category:    "one",
				Subcategory: "two",
			},
			valid: false,
		},
		{
			name: "invalid year",
			entry: Entry{
				Month:       11,
				Year:        2019,
				Category:    "one",
				Subcategory: "two",
			},
			valid: false,
		},
		{
			name: "invalid category",
			entry: Entry{
				Month:       11,
				Year:        -15,
				Category:    "",
				Subcategory: "two",
			},
			valid: false,
		},
		{
			name: "invalid subcategory",
			entry: Entry{
				Month:       11,
				Year:        -15,
				Category:    "one",
				Subcategory: "",
			},
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.entry.Validate() == nil) != tt.valid {
				t.Fatalf("expected to be valid: %t, got validation error: %v", tt.valid, tt.entry.Validate())
			}
		})
	}
}
