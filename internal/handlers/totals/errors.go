package totals

import public "github.com/Ozoniuss/casheer/pkg/casheerapi"

func NewGetRunningTotalFailedError(status int, detail string) public.Error {
	return public.Error{
		Title:  "Get Running Total Failed",
		Status: status,
		Detail: detail,
	}
}
