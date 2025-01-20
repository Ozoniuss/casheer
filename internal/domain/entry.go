package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/Ozoniuss/casheer/internal/domain/currency"
)

// Entry models an entry that can have zero or more expenses.
//
// An Entry, together with all Expenses, is an aggregate root, given that
// Entries and Expenses have to be strongly consistent. Changes to an Expense
// only happen via an Entry's public methods.
//
// Do not modify an entry directly. Use the provided functions to ensure
// entries are always valid.
type Entry struct {
	BaseModel
	currency.Value

	Month       Month
	Year        int
	Category    string
	Subcategory string
	Recurring   bool

	Expenses []Expense

	// Stores all changes that happened to this entry's expenses. One example
	// use case of this is database optimizations.
	ExpensesChanged []ExpenseChangedEvent
}

func NewEntry(month int, year int, category, subcategory string, recurring bool, value currency.Value) (Entry, error) {
	errs := make([]error, 0)

	mo, merr := NewMonth(month)
	if merr != nil {
		errs = append(errs, merr)
	}

	if category == "" {
		errs = append(errs, ErrEmptyEntryCategory)
	}
	if subcategory == "" {
		errs = append(errs, ErrEmptyEntrySubcategory)
	}

	if len(errs) != 0 {
		return Entry{}, ErrorInvalidEntry{underlying: slices.Clone(errs)}
	}

	return Entry{
		Month:       mo,
		Year:        year,
		Category:    category,
		Subcategory: subcategory,
		Recurring:   recurring,
		Value:       value,
		Expenses:    make([]Expense, 0),
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
		},
	}, nil
}

func (e *Entry) AddExpense(exp Expense) {
	e.Expenses = append(e.Expenses, exp)
	e.ExpensesChanged = append(e.ExpensesChanged, ExpenseChangedEvent{
		Data:   exp,
		Status: ExpenseCreated,
	})
}

func (e *Entry) ModifyExpense(exp Expense) error {
	for idx := range e.Expenses {
		if e.Expenses[idx].Id == exp.Id {
			e.Expenses[idx] = exp
			e.ExpensesChanged = append(e.ExpensesChanged, ExpenseChangedEvent{
				Data:   exp,
				Status: ExpenseModified,
			})
			return nil
		}
	}
	return ErrMissingExpense
}

func (e *Entry) DeleteExpense(exp Expense) error {
	pos := slices.IndexFunc(e.Expenses, func(expense Expense) bool {
		return expense.Id == exp.Id
	})
	if pos == -1 {
		return ErrMissingExpense
	}
	e.Expenses = slices.Delete(e.Expenses, pos, pos+1)
	return nil
}

type ExpenseChangedEvent struct {
	Data   Expense
	Status ExpenseStatus
}

type ExpenseStatus string

const (
	ExpenseCreated  ExpenseStatus = "created"
	ExpenseModified ExpenseStatus = "modified"
	ExpenseDeleted  ExpenseStatus = "deleted"
)

var ErrEmptyEntryCategory = errors.New("category must not be empty")
var ErrEmptyEntrySubcategory = errors.New("subcategory must not be empty")

var ErrMissingExpense = errors.New("the expense is not part of the aggregate")

type ErrorInvalidEntry = errorWithUnderlyingError
