package model

import (
	"gorm.io/gorm"
)

// Debt models a debt owed to or held by someone. Fields are automatically
// mapped by gorm to their database equivalents.
type Debt struct {
	BaseModel

	Person  string `validate:"required"`
	Amount  int    `validate:"required"`
	Details string
}

// AfterUpdate is a gorm hook that adds an error if the debt was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty debt.
func (d *Debt) AfterUpdate(tx *gorm.DB) (err error) {
	if d.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the debt was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty debt.
func (d *Debt) AfterDelete(tx *gorm.DB) (err error) {
	if d.Id == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}
