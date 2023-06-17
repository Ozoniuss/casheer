package errors

import (
	"fmt"
	"net/http"

	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func NewMissingParamError(missingParam string) public.Error {
	return public.Error{
		Title:  "Operation Failed",
		Status: http.StatusBadRequest,
		Detail: fmt.Sprintf("Missing URL parameter: %s", missingParam),
	}
}

func NewInvalidParamType(missingParam string) public.Error {
	return public.Error{
		Title:  "Operation Failed",
		Status: http.StatusBadRequest,
		Detail: fmt.Sprintf("Invalid type for URL parameter: %s", missingParam),
	}
}
