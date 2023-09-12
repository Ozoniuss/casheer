package calls

import (
	"net/http"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// MakeDELETE makes a simple DELETE request to the target url, and returns
// either a typed response or an error response.
func MakeDELETE[T casheerapi.DeleteDebtResponse](client *http.Client, url string) (T, error) {
	return makeRequest[T]("DELETE", client, url, nil)
}
