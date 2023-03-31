package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Expense models an expense that can be associated with an entry. Fields are
// automatically mapped by gorm to their database equivalents.
type Expense struct {
	BaseModel

	EntryId       uuid.UUID
	Value         float32 `validate:"required"`
	Name          string  `validate:"required"`
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
func RequiredEntry(expense *Expense) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var entries []Entry
		err := db.Session(&gorm.Session{}).Where("id = ?", expense.EntryId).Find(&entries).Error

		if err != nil {
			db.AddError(err)
		}

		if len(entries) == 0 {
			db.AddError(&NoEntryFoundErr{})
		}
		return db
	}
}
