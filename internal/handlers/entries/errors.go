package entries

import public "github.com/Ozoniuss/casheer/pkg/casheerapi"

func NewCreateEntryFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create Entry Failed",
		Status: status,
		Detail: detail,
	}
}

func NewDeleteEntryFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Delete Entry Failed",
		Status: status,
		Detail: detail,
	}
}

func NewListEntryFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "List Entry Failed",
		Status: status,
		Detail: detail,
	}
}

func NewUpdateEntryFailed(status int, detail string) public.Error {
	return public.Error{
		Title:  "Update Entry Failed",
		Status: status,
		Detail: detail,
	}
}

func NewGetEntryFailed(status int, detail string) public.Error {
	return public.Error{
		Title:  "Get Entry Failed",
		Status: status,
		Detail: detail,
	}
}
