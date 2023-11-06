package debts

import (
	"net/url"
	"strconv"

	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func DebtToPublic(debt model.Debt, debtsURL *url.URL) public.DebtData {
	return public.DebtData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(debt.Id),
			Type: public.DebtType,
		},
		Attributes: public.DebtAtrributes{
			Value: public.MonetaryValueAttributes{
				Amount:   debt.Amount,
				Exponent: debt.Exponent,
				Currency: debt.Currency,
			},
			Person:  debt.Person,
			Details: debt.Details,
			Timestamps: public.Timestamps{
				CreatedAt: debt.CreatedAt,
				UpdatedAt: debt.UpdatedAt,
			},
		},
		Links: public.DebtLinks{
			Self: debtsURL.JoinPath(strconv.Itoa(debt.Id)).String(),
		},
	}
}

func DebtToPublicList(debt model.Debt, debtsURL *url.URL) public.DebtListItemData {
	return public.DebtListItemData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(debt.Id),
			Type: public.DebtType,
		},
		Attributes: public.DebtAtrributes{
			Value: public.MonetaryValueAttributes{
				Amount:   debt.Amount,
				Exponent: debt.Exponent,
				Currency: debt.Currency,
			},
			Person:  debt.Person,
			Details: debt.Details,
			Timestamps: public.Timestamps{
				CreatedAt: debt.CreatedAt,
				UpdatedAt: debt.UpdatedAt,
			},
		},
		Links: public.DebtListItemLinks{
			Self: debtsURL.JoinPath(strconv.Itoa(debt.Id)).String(),
		},
	}
}
