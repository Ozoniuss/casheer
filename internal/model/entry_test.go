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
				Month:         10,
				Year:          2022,
				Category:      "one",
				Subcategory:   "two",
				ExpectedTotal: 100,
				RunningTotal:  0,
			},
			valid: true,
		},
		{
			id: 2,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:         13,
				Year:          2022,
				Category:      "one",
				Subcategory:   "two",
				ExpectedTotal: 100,
				RunningTotal:  0,
			},
			valid: false,
		},
		{
			id: 3,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:         11,
				Year:          2022,
				Category:      "",
				Subcategory:   "two",
				ExpectedTotal: 100,
				RunningTotal:  0,
			},
			valid: false,
		},
		{
			id: 4,
			entry: Entry{
				BaseModel: BaseModel{
					Id: 1,
				},
				Month:         11,
				Year:          -15,
				Category:      "one",
				Subcategory:   "two",
				ExpectedTotal: 100,
				RunningTotal:  0,
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
