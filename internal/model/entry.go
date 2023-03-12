package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entry models an entry of a planning. Fields are automatically mapped by gorm
// to their database equivalents.
type Entry struct {
	BaseModel

	// Postgresql doesn't support unsigned int.
	Month         int8   `validate:"required,gte=1,lte=12"`
	Year          int16  `validate:"required,gte=2020"`
	Category      string `validate:"required"`
	Subcategory   string `validate:"required"`
	ExpectedTotal float32
	RunningTotal  float32
	Recurring     bool
}

// AfterUpdate is a gorm hook that adds an error if the entry was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty entry.
func (e *Entry) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the entry was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty entry.
func (e *Entry) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}
