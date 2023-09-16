package calls

import (
	"net/http"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// MakePOST makes a simple POST request to the target url, and returns
// either a typed response or an error response.
func MakePOST[T casheerapi.CreateDebtResponse | casheerapi.CreateEntryResponse](client *http.Client, url string, payload []byte) (T, error) {
	return makeRequest[T]("POST", client, url, payload)
}
