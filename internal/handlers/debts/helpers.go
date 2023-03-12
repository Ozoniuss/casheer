package debts

import (
	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func DebtToPublic(debt model.Debt) public.DebtData {
	return public.DebtData{
		ResourceID: public.ResourceID{
			Id:   debt.Id.String(),
			Type: public.DebtType,
		},
		Person:  debt.Person,
		Amount:  debt.Amount,
		Details: debt.Details,
		Timestamps: public.Timestamps{
			CreatedAt: debt.CreatedAt,
			UpdatedAt: debt.UpdatedAt,
		},
	}
}
