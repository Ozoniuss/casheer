package entries

import public "github.com/Ozoniuss/casheer/pkg/casheerapi"

func NewCreateEntryFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Create Entry Failed",
		Status: status,
		Detail: detail,
	}
}
