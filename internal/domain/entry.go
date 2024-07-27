package domain

import (
	"errors"
	"time"

	"github.com/Ozoniuss/casheer/internal/domain/currency"
)

// Entry models an entry that can have zero or more expenses.
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
}

func NewEntry(month int, year int, category, subcategory string, recurring bool) (Entry, error) {
	var err error = nil

	mo, merr := NewMonth(month)
	if merr != nil {
		errors.Join(err, merr)
	}

	if category == "" {
		err = errors.Join(err, ErrEmptyCategory)
	}
	if subcategory == "" {
		err = errors.Join(err, ErrEmptySubcategory)
	}

	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Month:       mo,
		Year:        year,
		Category:    category,
		Subcategory: subcategory,
		Recurring:   recurring,
		Expenses:    make([]Expense, 0),
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
		},
	}, nil
}

func (e *Entry) AddExpense(exp Expense) {
	e.Expenses = append(e.Expenses, exp)
}

var ErrEmptyCategory = errors.New("category must not be empty")
var ErrEmptySubcategory = errors.New("subcategory must not be empty")
