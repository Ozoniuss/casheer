package model

import (
	"fmt"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/go-playground/validator/v10"
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

// ValidateModelScope can be used as a scope to validate any of the existing
// GORM models. It makes use of go-validator.
//
// Will be fully deprecated once the new validation works for all models.
func ValidateModelScope[T Entry | Expense | Debt](obj T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v := validator.New()
		if err := v.Struct(obj); err != nil {
			var dberr error
			switch o := any(obj).(type) {
			case Entry:
				dberr = ierrors.NewInvalidModelError("entry", err)
			case Expense:
				dberr = ierrors.NewInvalidModelError("expense", err)
			case Debt:
				fmt.Println("ce dq ce drq ce drqs")
				dberr = o.Validate()
				fmt.Println(dberr)
			default:
				panic("this should not happen but it's not handled yet")
			}
			db.AddError(dberr)
		}
		return db
	}
}
