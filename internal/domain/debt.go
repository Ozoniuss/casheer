package domain

import (
	"errors"
	"time"

	"github.com/Ozoniuss/casheer/internal/domain/currency"
)

// Debt models a debt owed to or held by someone. Fields are automatically
// mapped by gorm to their database equivalents.
type Debt struct {
	BaseModel
	currency.Value

	Person  string
	Details string
}

func NewDebt(person string, details string, value currency.Value) (Debt, error) {
	errs := make([]error, 0)

	if person == "" {
		errs = append(errs, ErrMissingDebtPerson)
	}

	if len(errs) != 0 {
		return Debt{}, ErrInvalidDebt{underlying: errs}
	}

	return Debt{
		Person:  person,
		Details: details,
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
		},
		Value: value,
	}, nil
}

var ErrMissingDebtPerson = errors.New("person must not be empty")

type ErrInvalidDebt = errorWithUnderlyingError
