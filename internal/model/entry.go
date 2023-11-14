package model

import (
	"gorm.io/gorm"
)

// Entry models an entry of a planning. Fields are automatically mapped by gorm
// to their database equivalents.
type Entry struct {
	BaseModel
	Value

	Month       int
	Year        int
	Category    string
	Subcategory string
	Recurring   bool

	Expenses []Expense
}

func (e Entry) Validate() error {
	b := NewBaseModelErrorBuilder("entry")
	if e.Month > 12 || e.Month < 1 {
		b.AddError("month must be between 1 and 12")
	}
	if e.Year < 2020 {
		b.AddError("year must be at least 2020")
	}
	if len(e.Category) == 0 {
		b.AddError("category cannot be empty")
	}
	if len(e.Subcategory) == 0 {
		b.AddError("subcategory cannot be empty")
	}
	return b.Error()
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
