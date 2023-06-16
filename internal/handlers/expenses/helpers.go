package expenses

import (
	"fmt"
	"strconv"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// ExpenseToPublic converts an expense to its public API representation.
func ExpenseToPublic(expense model.Expense, apiPaths config.ApiPaths) public.ExpenseData {
	return public.ExpenseData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(expense.Id),
			Type: public.EntryType,
		},
		Attributes: public.ExpenseAttributes{
			Value:         expense.Value,
			Name:          expense.Name,
			Description:   expense.Description,
			PaymentMethod: expense.PaymentMethod,
			Timestamps: public.Timestamps{
				CreatedAt: expense.CreatedAt,
				UpdatedAt: expense.UpdatedAt,
			},
		},
		Links: public.ExpenseLinks{
			Self:       fmt.Sprintf("%s/%s/expenses/%s", apiPaths.Entries, expense.EntryId.String(), strconv.Itoa(expense.Id)),
			Collection: fmt.Sprintf("%s/%s/expenses/", apiPaths.Entries, expense.EntryId.String()),
		},
	}
}

// ExpenseToPublic converts an expense to its public API representation in
// expense listings.
func ExpenseToPublicList(expense model.Expense, apiPaths config.ApiPaths) public.ExpenseListItemData {
	return public.ExpenseListItemData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(expense.Id),
			Type: public.EntryType,
		},
		Attributes: public.ExpenseAttributes{
			Value:         expense.Value,
			Name:          expense.Name,
			Description:   expense.Description,
			PaymentMethod: expense.PaymentMethod,
			Timestamps: public.Timestamps{
				CreatedAt: expense.CreatedAt,
				UpdatedAt: expense.UpdatedAt,
			},
		},
		Links: public.ExpenseListItemLinks{
			Self: fmt.Sprintf("%s/%s/expenses/%s", apiPaths.Entries, expense.EntryId.String(), strconv.Itoa(expense.Id)),
		},
	}
}
