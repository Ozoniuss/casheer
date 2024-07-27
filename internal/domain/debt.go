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

func NewDebt(person string, details string) (Debt, error) {

	if person == "" {
		return Debt{}, ErrEmptyPerson
	}

	return Debt{
		Person:  person,
		Details: details,
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
		},
	}, nil
}

var ErrEmptyPerson = errors.New("person must not be empty")
