package errors

import (
	"fmt"
	"net/http"
	"strconv"

	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func NewInvalidURLError(detail string) public.Error {
	return public.Error{
		Title:  "Invalid URL",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func NewRequestBindingError(detail string) public.Error {
	return public.Error{
		Title:  "Invalid JSON Body",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func NewQueryParamsBindingError(detail string) public.Error {
	return public.Error{
		Title:  "Invalid Query Params",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func NewMissingContextParamError(detail string) public.Error {
	return public.Error{
		Title:  "Missing Context Parameter",
		Status: http.StatusInternalServerError,
		Detail: detail,
	}
}

func NewInvalidContextParamError(detail string) public.Error {
	return public.Error{
		Title:  "Invalid Context Parameter",
		Status: http.StatusInternalServerError,
		Detail: detail,
	}
}

func NewMissingParamError(missingParam string) public.Error {
	return public.Error{
		Title:  "Missing URL Parameter",
		Status: http.StatusBadRequest,
		Detail: fmt.Sprintf("Missing URL parameter: %s", strconv.QuoteToASCII(missingParam)),
	}
}

func NewInvalidParamType(paramName string) public.Error {
	return public.Error{
		Title:  "Invalid URL Parameter",
		Status: http.StatusBadRequest,
		Detail: fmt.Sprintf("URL parameter %s is not an integer.", strconv.QuoteToASCII(paramName)),
	}
}

func NewUnknownError(detail string) public.Error {
	return public.Error{
		Title:  "Unknown Error",
		Status: http.StatusInternalServerError,
		Detail: fmt.Sprintf("An unknown error occured: %s", detail),
	}
}

func NewNotFoundError[T string | int](resource T) public.Error {
	return public.Error{
		Title:  "Not Found Error",
		Status: http.StatusNotFound,
		Detail: fmt.Sprintf("Resource %v was not found.", resource),
	}
}
