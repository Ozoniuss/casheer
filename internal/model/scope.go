package model

import (
	"gorm.io/gorm"
)

// ValidateModel can be used as a scope to validate any of the existing GORM
// models, that implement the ModelValidator interface.
func ValidateModel(model ModelValidator) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.AddError(model.Validate())
		return db
	}
}
