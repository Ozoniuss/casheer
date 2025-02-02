package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/Ozoniuss/casheer/internal/domain/currency"
)

// Expense models an expense that can be associated with an entry.
//
// Do not modify an expense directly. Use the provided functions to ensure
// expenses are always valid.
type Expense struct {
	BaseModel
	currency.Value

	Name          string
	Description   string
	PaymentMethod string
}

func NewScratchExpense(name string, description string, paymentMethod string, value currency.Value) (Expense, error) {
	errs := make([]error, 0)
	if name == "" {
		errs = append(errs, ErrEmptyExpenseName)
	}
	if len(errs) != 0 {
		return Expense{}, ErrorInvalidExpense{underlying: slices.Clone(errs)}
	}
	return Expense{
		Name:          name,
		Description:   description,
		PaymentMethod: paymentMethod,
		Value:         value,
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
		},
	}, nil
}

var ErrEmptyExpenseName = errors.New("expense must have a name")

type ErrorInvalidExpense = errorWithUnderlyingError
