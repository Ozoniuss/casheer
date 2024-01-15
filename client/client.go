package client

import (
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// CasheerClient is a client provided for casheer which can be used to abstract
// away the process of creating the HTTP requests towards the API. It also
// provides some methods that are more oriented towards the bussiness
// requirements, as opposed to the generality provided by the public package.
type CasheerClient interface {
	CreateDebt(person string, details string, amount string, currency string, exponent string) (public.CreateDebtResponse, error)
	DeleteDebt(debtId int) public.DeleteDebtResponse
	UpdateDebt(debtId int, person *string, details *string, amount *string, currency *string, exponent *string) (public.UpdateDebtResponse, error)
	GetDebt(debtId int) (public.GetDebtResponse, error)
	ListDebts() (public.ListDebtResponse, error)
	ListDebtsForPerson(person string) (public.ListDebtResponse, error)

	CreateEntryWithCurrency(month int, year int, category string, subcategory string, expected_total int, currency string, recurring bool) (public.CreateEntryResponse, error)
	ListEntries() (public.ListEntryResponse, error)
}
