package debts

import (
	"fmt"
	"net/http"

	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func NewInvalidDebtError(reason string) public.Error {
	return public.Error{
		Title:  "Invalid Debt",
		Status: http.StatusBadRequest,
		Detail: fmt.Sprintf("Could not validate debt: %s", reason),
	}
}

func NewCreateDebtFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create Debt Failed",
		Status: status,
		Detail: detail,
	}
}

func NewDeleteDebtFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Delete Debt Failed",
		Status: status,
		Detail: detail,
	}
}

func NewListDebtFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "List Debt Failed",
		Status: status,
		Detail: detail,
	}
}

func NewUpdateDebtFailed(status int, detail string) public.Error {
	return public.Error{
		Title:  "Update Debt Failed",
		Status: status,
		Detail: detail,
	}
}

func NewGetDebtFailed(status int, detail string) public.Error {
	return public.Error{
		Title:  "Get Debt Failed",
		Status: status,
		Detail: detail,
	}
}
