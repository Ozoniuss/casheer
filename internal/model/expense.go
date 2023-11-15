package model

import (
	"fmt"

	"gorm.io/gorm"
)

// Expense models an expense that can be associated with an entry. Fields are
// automatically mapped by gorm to their database equivalents.
type Expense struct {
	BaseModel
	Value

	EntryId       int    `validate:"required"`
	Name          string `validate:"required"`
	Description   string
	PaymentMethod string
}

func (e Expense) Validate() error {
	b := NewBaseModelErrorBuilder("entry")
	if len(e.Name) == 0 {
		b.AddError("name cannot be empty")
	}
	return b.Error()
}

type ErrExpenseInvalidEntryKey struct {
	entryKey int
}

func (e ErrExpenseInvalidEntryKey) Error() string {
	return fmt.Sprintf("entry with id %d was not found", e.entryKey)
}

// RequiredEntry can be used as a scope to return a custom error if the
// entry id associated with the expense doesn't exist.
func RequiredEntry(entryId int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var entries []Entry
		err := db.Session(&gorm.Session{NewDB: true}).Table("entries").Where("id = ?", entryId).Find(&entries).Error

		if err != nil {
			db.AddError(err)
		}

		if len(entries) == 0 {
			db.AddError(ErrExpenseInvalidEntryKey{
				entryKey: entryId,
			})
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
