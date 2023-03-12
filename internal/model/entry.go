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
	Month         int8
	Year          int16
	Category      string
	Subcategory   string
	ExpectedTotal float32
	RunningTotal  float32
	Recurring     bool
}

func (e *Entry) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}
