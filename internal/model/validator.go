package model

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ValidateModel can be used as a scope to validate any of the existing GORM
// models. It makes use of go-validator.
func ValidateModel[T Entry](obj T, adderr error) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v := validator.New()
		if err := v.Struct(obj); err != nil {
			db.AddError(adderr)
		}
		return db
	}
}

// ValidateModelFields does the same as ValidateModel, except it only validates
// the provided entity fields. Note that the provided fields should be
// uppercase.
func ValidateModelFields[T Entry](obj T, fields []string, adderr error) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v := validator.New()
		if err := v.StructPartial(obj, fields...); err != nil {
			db.AddError(InvalidEntryErr{})
		}
		return db
	}
}
