package model

import (
	"github.com/Ozoniuss/casheer/internal/currency"
	"gorm.io/gorm"
)

// Expense models an expense that can be associated with an entry. Fields are
// automatically mapped by gorm to their database equivalents.
type Expense struct {
	BaseModel

	EntryId int `validate:"required"`
	currency.Value
	Name          string `validate:"required"`
	Description   string
	PaymentMethod string
}

type NoEntryFoundErr struct {
}

func (e *NoEntryFoundErr) Error() string {
	return "no entry found"
}

// RequiredEntry can be used as a scope to return a custom error if the
// entry id associated with the expense doesn't exist.
func RequiredEntry(entryId int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var entries []Entry
		err := db.Session(&gorm.Session{}).Where("id = ?", entryId).Find(&entries).Error

		if err != nil {
			db.AddError(err)
		}

		if len(entries) == 0 {
			db.AddError(&NoEntryFoundErr{})
		}
		return db
	}
}

// AfterUpdate is a gorm hook that adds an error if the expense was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty expense.
func (e *Expense) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the expense was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty expense.
func (e *Expense) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}
