package expenses

import (
	"net/url"
	"strconv"

	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// ExpenseToPublic converts an expense to its public API representation.
func ExpenseToPublic(expense model.Expense, entriesURL *url.URL) public.ExpenseData {
	return public.ExpenseData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(expense.Id),
			Type: public.EntryType,
		},
		Attributes: public.ExpenseAttributes{
			Amount:        expense.Amount,
			Currency:      expense.Currency,
			Exponent:      expense.Exponent,
			Name:          expense.Name,
			Description:   expense.Description,
			PaymentMethod: expense.PaymentMethod,
			Timestamps: public.Timestamps{
				CreatedAt: expense.CreatedAt,
				UpdatedAt: expense.UpdatedAt,
			},
		},
		Links: public.ExpenseLinks{
			Self:       entriesURL.JoinPath(strconv.Itoa(expense.EntryId), "expenses", strconv.Itoa(expense.Id)).String(),
			Collection: entriesURL.JoinPath(strconv.Itoa(expense.EntryId), "expenses/").String(),
		},
	}
}

// ExpenseToPublic converts an expense to its public API representation in
// expense listings.
func ExpenseToPublicList(expense model.Expense, entriesURL *url.URL) public.ExpenseListItemData {
	return public.ExpenseListItemData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(expense.Id),
			Type: public.EntryType,
		},
		Attributes: public.ExpenseAttributes{
			Amount:        expense.Amount,
			Currency:      expense.Currency,
			Exponent:      expense.Exponent,
			Name:          expense.Name,
			Description:   expense.Description,
			PaymentMethod: expense.PaymentMethod,
			Timestamps: public.Timestamps{
				CreatedAt: expense.CreatedAt,
				UpdatedAt: expense.UpdatedAt,
			},
		},
		Links: public.ExpenseListItemLinks{
			Self: entriesURL.JoinPath(strconv.Itoa(expense.EntryId), "expenses", strconv.Itoa(expense.Id)).String(),
		},
	}
}
