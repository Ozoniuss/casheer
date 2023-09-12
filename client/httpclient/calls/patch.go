package calls

import (
	"net/http"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// MakePATCH makes a simple PATCH request to the target url, and returns
// either a typed response or an error response.
func MakePATCH[T casheerapi.UpdateDebtResponse](client *http.Client, url string, payload []byte) (T, error) {
	return makeRequest[T]("PATCH", client, url, payload)
}
