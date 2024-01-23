package entries

import (
	"net/url"
	"strconv"

	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// computeRunningTotal returns the running total of a list of expenses.
func computeRunningTotal(expenses []model.Expense) int {
	var rt int = 0
	for _, exp := range expenses {
		// TODO: take into account currency, if doing so
		rt += exp.Value.Amount
	}
	return rt
}

func EntryToPublic(entry model.Entry, entriesURL *url.URL, runningTotal int) public.EntryData {
	return public.EntryData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(entry.Id),
			Type: public.EntryType,
		},
		Attributes: public.EntryAtrributes{
			Year:        int(entry.Year),
			Month:       int(entry.Month),
			Category:    entry.Category,
			Subcategory: entry.Subcategory,
			ExpectedTotal: public.MonetaryValueAttributes{
				Currency: entry.Currency,
				Amount:   entry.Amount,
				Exponent: entry.Exponent,
			},
			Recurring: entry.Recurring,
			Timestamps: public.Timestamps{
				CreatedAt: entry.CreatedAt,
				UpdatedAt: entry.UpdatedAt,
			},
		},
		Meta: public.EntryMeta{
			RunningTotal: runningTotal,
		},
		Links: public.EntryLinks{
			Self: entriesURL.JoinPath(strconv.Itoa(entry.Id)).String(),
		},
		Relationships: public.EntryRelationships{
			Expenses: public.EntryExpenseRelationship{
				Links: public.EntryExpenseRelationshipLinks{
					Related: entriesURL.JoinPath(strconv.Itoa(entry.Id), "expenses/").String(),
				},
			},
		},
	}
}

func IncludedExpensesToPublic(expense model.Expense, entriesURL *url.URL) public.IncludedExpenseData {
	return public.IncludedExpenseData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(expense.Id),
			Type: public.ExpenseType,
		},
		Attributes: public.ExpenseAttributes{
			Value: public.MonetaryValueAttributes{
				Amount:   expense.Amount,
				Currency: expense.Currency,
				Exponent: expense.Exponent,
			},
			Name:          expense.Name,
			Description:   expense.Description,
			PaymentMethod: expense.PaymentMethod,
			Timestamps: public.Timestamps{
				CreatedAt: expense.CreatedAt,
				UpdatedAt: expense.UpdatedAt,
			},
		},
		Links: public.ExpenseLinks{
			Self: entriesURL.JoinPath(strconv.Itoa(expense.EntryId), "expenses", strconv.Itoa(expense.Id)).String(),
		},
	}
}

func EntryToPublicList(entry model.Entry, entriesURL *url.URL, runningTotal int) public.EntryListItemData {
	return public.EntryListItemData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(entry.Id),
			Type: public.EntryType,
		},
		Attributes: public.EntryAtrributes{
			Year:        int(entry.Year),
			Month:       int(entry.Month),
			Category:    entry.Category,
			Subcategory: entry.Subcategory,
			ExpectedTotal: public.MonetaryValueAttributes{
				Currency: entry.Currency,
				Amount:   entry.Amount,
				Exponent: entry.Exponent,
			},
			Recurring: entry.Recurring,
			Timestamps: public.Timestamps{
				CreatedAt: entry.CreatedAt,
				UpdatedAt: entry.UpdatedAt,
			},
		},
		Meta: public.EntryMeta{
			RunningTotal: runningTotal,
		},
		Links: public.EntryListItemLinks{
			Self: entriesURL.JoinPath(strconv.Itoa(entry.Id)).String(),
		},
		Relationships: public.EntryRelationships{
			Expenses: public.EntryExpenseRelationship{
				Links: public.EntryExpenseRelationshipLinks{
					Related: entriesURL.JoinPath(strconv.Itoa(entry.Id), "expenses/").String(),
				},
			},
		},
	}
}
