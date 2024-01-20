package store

import (
	"github.com/Ozoniuss/casheer/internal/domain/model"
)

type Store interface {
	AddEntry(model.Entry) (model.Entry, error)
	DeleteEntry(model.Id) (model.Entry, error)
	ListEntries() ([]model.Entry, error)
	GetEntry(model.Id) (model.Entry, error)
	UpdateEntry(model.Id) (model.Entry, error)

	AddExpenseToEntry(model.Id, model.Expense) (model.Expense, error)
	DeleteExpense(model.Id) (model.Expense, error)
	ListExpensesForEntry(model.Id) ([]model.Expense, error)
	GetExpense(model.Id) (model.Expense, error)
	UpdateExpense(model.Id) (model.Expense, error)

	AddDebt(model.Debt) (model.Debt, error)
	DeleteDebt(model.Id) (model.Debt, error)
	ListDebts() ([]model.Debt, error)
	GetDebt(model.Id) (model.Debt, error)
	UpdateDebt(model.Id) (model.Debt, error)
}

type ErrAlreadyExists struct {
}

func (err ErrAlreadyExists) Error() string {
	return ""
}

type ErrNotFound struct {
}

func (err ErrNotFound) Error() string {
	return ""
}

type ErrInvalidConstrains struct {
}

func (err ErrInvalidConstrains) Error() string {
	return ""
}

type ErrUnknown struct {
}

func (err ErrUnknown) Error() string {
	return ""
}
