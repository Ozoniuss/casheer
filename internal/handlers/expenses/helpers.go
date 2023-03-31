package expenses

import (
	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// ExpenseToPublic converts an expense to its public API representation.
func ExpenseToPublic(expense model.Expense) public.ExpenseData {
	return public.ExpenseData{
		ResourceID: public.ResourceID{
			Id:   expense.Id.String(),
			Type: public.EntryType,
		},
		Value:         expense.Value,
		Name:          expense.Name,
		Description:   expense.Description,
		PaymentMethod: expense.PaymentMethod,
		Timestamps: public.Timestamps{
			CreatedAt: expense.CreatedAt,
			UpdatedAt: expense.UpdatedAt,
		},
	}
}
