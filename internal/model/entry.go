package model

import (
	"gorm.io/gorm"
)

// Entry models an entry of a planning. Fields are automatically mapped by gorm
// to their database equivalents.
//
// TODO: custom validation message.
type Entry struct {
	BaseModel

	// Postgresql doesn't support unsigned int.
	Month         int    `validate:"required,gte=1,lte=12"`
	Year          int    `validate:"required,gte=2020"`
	Category      string `validate:"required" json:"category"`
	Subcategory   string `validate:"required"`
	ExpectedTotal int
	Recurring     bool

	Expenses []Expense
}

type InvalidEntryErr struct {
}

func (e InvalidEntryErr) Error() string {
	return "invalid entry"
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
