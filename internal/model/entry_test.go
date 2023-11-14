package model

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidEntry(t *testing.T) {
	tests := []struct {
		id    int
		entry Entry
		valid bool
	}{
		{
			id: 1,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:       10,
				Year:        2022,
				Category:    "one",
				Subcategory: "two",
				Value: Value{
					Currency: "EUR",
					Amount:   100,
					Exponent: 0,
				},
			},
			valid: true,
		},
		{
			id: 2,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:       13,
				Year:        2022,
				Category:    "one",
				Subcategory: "two",
				Value: Value{
					Currency: "EUR",
					Amount:   100,
					Exponent: 0,
				},
			},
			valid: false,
		},
		{
			id: 3,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:       11,
				Year:        2022,
				Category:    "",
				Subcategory: "two",
				Value: Value{
					Currency: "EUR",
					Amount:   100,
					Exponent: 0,
				},
			},
			valid: false,
		},
		{
			id: 4,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:       11,
				Year:        -15,
				Category:    "one",
				Subcategory: "two",
				Value: Value{
					Currency: "EUR",
					Amount:   100,
					Exponent: 0,
				},
			},
			valid: false,
		},
	}
	validate := validator.New()
	for _, tt := range tests {
		err := validate.Struct(tt.entry)
		if tt.valid {
			if err != nil {
				t.Fatalf("test %d failed: got unwanted validation error for entry: %s", tt.id, err.Error())
			}
		} else {
			if err == nil {
				t.Fatalf("test %d failed: entry should not be valid", tt.id)
			}
		}
	}
}

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
