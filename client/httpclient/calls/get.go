package calls

import (
	"net/http"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// MakeGET makes a simple GET request to the target url, and returns either a
// typed response or an error response.
func MakeGET[T casheerapi.GetDebtResponse](client *http.Client, url string) (T, error) {
	return makeRequest[T]("GET", client, url, nil)
}
