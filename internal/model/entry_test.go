package model

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// genUuidPattern generates an uuid based on the provided pattern. It is assumed
// that the pattern has lenght 1,2 or 4.
func genUuidPattern(pattern string) uuid.UUID {
	if len(pattern) != 1 && len(pattern) != 2 && len(pattern) != 4 {
		panic(fmt.Sprintf("invalid pattern length %d", len(pattern)))
	}
	if len(pattern) == 1 {
		pattern = fmt.Sprintf("%[1]s%[1]s%[1]s%[1]s", pattern)
	} else if len(pattern) == 2 {
		pattern = fmt.Sprintf("%[1]s%[1]s", pattern)
	}
	return uuid.MustParse(fmt.Sprintf("%[1]s%[1]s-%[1]s-%[1]s-%[1]s-%[1]s%[1]s%[1]s", pattern))
}

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
