package expenses

import (
	"net/http"

	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func NewCreateExpenseFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create Expense Failed",
		Status: status,
		Detail: detail,
	}
}

func NewDeleteExpenseFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Delete Expense Failed",
		Status: status,
		Detail: detail,
	}
}

func NewListExpenseFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "List Expense Failed",
		Status: status,
		Detail: detail,
	}
}

func NewUpdateExpenseFailed(status int, detail string) public.Error {
	return public.Error{
		Title:  "Update Expense Failed",
		Status: status,
		Detail: detail,
	}
}

func NewGetExpenseFailed(status int, detail string) public.Error {
	return public.Error{
		Title:  "Get Expense Failed",
		Status: status,
		Detail: detail,
	}
}

func NewInvalidEntryError(detail string) public.Error {
	return public.Error{
		Title:  "Operation Failed",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}
