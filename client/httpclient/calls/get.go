package calls

import (
	"net/http"
)

// MakeGET makes a simple GET request to the target url, and returns either a
// typed response or an error response.
func MakeGET[T ResponseType](client *http.Client, url string) (T, error) {
	return makeRequest[T]("GET", client, url, nil, nil)
}
