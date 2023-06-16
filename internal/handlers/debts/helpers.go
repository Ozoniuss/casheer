package debts

import (
	"fmt"
	"strconv"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func DebtToPublic(debt model.Debt, paths config.ApiPaths) public.DebtData {
	return public.DebtData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(debt.Id),
			Type: public.DebtType,
		},
		Attributes: public.DebtAtrributes{
			Person:  debt.Person,
			Amount:  debt.Amount,
			Details: debt.Details,
			Timestamps: public.Timestamps{
				CreatedAt: debt.CreatedAt,
				UpdatedAt: debt.UpdatedAt,
			},
		},
		Links: public.DebtLinks{
			Self:       fmt.Sprintf("%s/%d", paths.Debts, debt.Id),
			Collection: paths.Debts + "/",
		},
	}
}

func DebtToPublicList(debt model.Debt, paths config.ApiPaths) public.DebtListItemData {
	return public.DebtListItemData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(debt.Id),
			Type: public.DebtType,
		},
		Attributes: public.DebtAtrributes{
			Person:  debt.Person,
			Amount:  debt.Amount,
			Details: debt.Details,
			Timestamps: public.Timestamps{
				CreatedAt: debt.CreatedAt,
				UpdatedAt: debt.UpdatedAt,
			},
		},
		Links: public.DebtListItemLinks{
			Self: fmt.Sprintf("%s/%d", paths.Debts, debt.Id),
		},
	}
}
