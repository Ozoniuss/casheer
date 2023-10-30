package model

import (
	"strings"

	"gorm.io/gorm"
)

// Entry models an entry of a planning. Fields are automatically mapped by gorm
// to their database equivalents.
//
// TODO: custom validation message.
type Entry struct {
	BaseModel
	Value

	Month       int    `validate:"required,gte=1,lte=12"`
	Year        int    `validate:"required,gte=2020"`
	Category    string `validate:"required" json:"category"`
	Subcategory string `validate:"required"`
	Recurring   bool

	Expenses []Expense
}

func ValidateEntry(e Entry) error {
	reasons := []string{}
	if e.Month < 2020 {
		reasons = append(reasons, "year must be after 2020")
	}
	if e.Month < 1 || e.Month > 12 {
		reasons = append(reasons, "month must be between 1 and 12")
	}
	if e.Category == "" {
		reasons = append(reasons, "category must not be empty")
	}
	if e.Subcategory == "" {
		reasons = append(reasons, "subcategory must not be empty")
	}
	return NewInvalidEntryErr(reasons)
}

type InvalidEntryErr struct {
	reasons []string
}

func NewInvalidEntryErr(reasons []string) error {
	if len(reasons) == 0 {
		return nil
	}
	return InvalidEntryErr{
		reasons: reasons,
	}
}

func (e InvalidEntryErr) Error() string {
	b := &strings.Builder{}
	b.WriteString("invalid entry: ")
	for idx, reason := range e.reasons {
		b.WriteString(reason)
		if idx != len(e.reasons)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

// AfterUpdate is a gorm hook that adds an error if the entry was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty entry.
func (e *Entry) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the entry was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty entry.
func (e *Entry) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}
