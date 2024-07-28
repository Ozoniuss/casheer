package store

import (
	"context"

	"github.com/Ozoniuss/casheer/internal/domain"
)

// CasheerStore models the interaction required with a casheer database.
type CasheerStore interface {
	SaveDebt(context.Context, domain.Debt) (domain.Debt, error)
	LoadDebt(context.Context, int) (domain.Debt, error)
	ListDebts(context.Context) ([]domain.Debt, error)

	SaveEntry(context.Context, domain.Entry) (domain.Entry, error)
	LoadEntry(context.Context, int) (domain.Entry, error)
	ListEntries(context.Context) ([]domain.Entry, error)
}

type ErrNotFound struct {
	Details string
	Orig    error
}

func (e ErrNotFound) Error() string {
	return e.Details
}
func (e ErrNotFound) Unwrap() error {
	return e.Orig
}
