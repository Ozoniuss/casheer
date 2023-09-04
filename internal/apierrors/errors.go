package apierrors

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

func NewInvalidURLNoTrailingSlashError() public.Error {
	return public.Error{
		Title:  "Invalid URL",
		Status: http.StatusBadRequest,
		Detail: "Please add a trailing slash.",
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

func NewNotFoundError() public.Error {
	return public.Error{
		Title:  "Resource Not Found",
		Status: http.StatusNotFound,
		Detail: "Resource was not found.",
	}
}

func NewInvalidResourceError(detail string) public.Error {
	return public.Error{
		Title:  "Invalid Resource",
		Status: http.StatusUnprocessableEntity,
		Detail: detail,
	}
}

func NewAlreadyExistsError(detail string) public.Error {
	return public.Error{
		Title:  "Resource Already Exists",
		Status: http.StatusConflict,
		Detail: detail,
	}
}
