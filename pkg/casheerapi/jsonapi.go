package casheerapi

import (
	"fmt"
	"strconv"
)

type ResourceID struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

// Error represents an error message that is sent to the client in case the
// HTTP request returns an error.
//
// The error message is modelled after the jsonapi specification:
// https://jsonapi.org/format/#error-objects. Many of the suggested fields had
// not been included, since they were unnecessary for the size of this project.
// In production, however, it is generally a good idea to stick to a
// well-defined standard.
//
// See the Error struct below for the actual error content.
type ErrorResponse struct {
	Err Error `json:"error"`
}

// Error implements the error interface for error responses. It allows go
// clients to treat error responses as actual errors, if desired.
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("{error:%v}", e.Err)
}

// Unwrap is useful for go API clients, allowing to retrieve the underlying
// public error directly from the error response.
func (e ErrorResponse) Unwrap() error {
	return e.Err
}

type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Error implements the error interface for public errors. This has two
// benefits:
//
// - Go clients of the casheer API can unwrap the error from an ErrorResponse
// to an actual error;
// - In handlers that contain operations such as transactions, a public error
// can be set directly inside the transaction. Callers can check if that
// operation return a public error directly, instead of having to build one
// based on some internal errors encountered in the transaction.
func (e Error) Error() string {
	return fmt.Sprintf("{\"Status\":%d,\"Title\":%s,\"Detail\":%s}",
		e.Status, strconv.Quote(e.Title), strconv.Quote(e.Detail))
}
