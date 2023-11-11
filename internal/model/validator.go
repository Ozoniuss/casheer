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
			case Expense:
				dberr = ierrors.NewInvalidModelError("expense", err)
			case Debt:
				dberr = ierrors.NewInvalidModelError("debt", err)
			default:
				panic("this should not happen but it's not handled yet")
			}
			db.AddError(dberr)
		}
		return db
	}
}
