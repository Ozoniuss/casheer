package model

import (
	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ValidateModelScope can be used as a scope to validate any of the existing
// GORM models. It makes use of go-validator.
func ValidateModelScope[T Entry | Expense | Debt](obj T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v := validator.New()
		if err := v.Struct(obj); err != nil {
			var dberr error
			switch any(obj).(type) {
			case Entry:
				dberr = ierrors.NewInvalidModelError("entry", err)
				break
			case Expense:
				dberr = ierrors.NewInvalidModelError("expense", err)
				break
			case Debt:
				dberr = ierrors.NewInvalidModelError("debt", err)
				break
			default:
				panic("this should not happen but it's not handled yet")
			}
			db.AddError(dberr)
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
