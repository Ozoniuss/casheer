package httpclient

import (
	"encoding/json"
	"fmt"

	"github.com/Ozoniuss/casheer/client/httpclient/calls"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (c *CasheerHTTPClient) CreateEntry(month int, year int, category string, subcategory string, recurring bool) (public.CreateEntryResponse, error) {
	requestURL := c.debtsURL.String()

	req := public.CreateEntryRequest{
		Data: struct {
			Type       string "json:\"type\" binding:\"required\""
			Attributes struct {
				Month         *int   "json:\"month,omitempty\""
				Year          *int   "json:\"year,omitempty\""
				Category      string "json:\"category\" binding:\"required\""
				Subcategory   string "json:\"subcategory\" binding:\"required\""
				ExpectedTotal int    "json:\"expected_total\" binding:\"required\""
				Recurring     bool   "json:\"recurring\""
			} "json:\"attributes\""
		}{
			Type: "entry",
			Attributes: struct {
				Month         *int   "json:\"month,omitempty\""
				Year          *int   "json:\"year,omitempty\""
				Category      string "json:\"category\" binding:\"required\""
				Subcategory   string "json:\"subcategory\" binding:\"required\""
				ExpectedTotal int    "json:\"expected_total\" binding:\"required\""
				Recurring     bool   "json:\"recurring\""
			}{
				Month:       &month,
				Year:        &year,
				Category:    category,
				Subcategory: subcategory,
				Recurring:   recurring,
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return public.CreateEntryResponse{}, fmt.Errorf("marshalling request into JSON: %s", err.Error())
	}

	return calls.MakePOST[public.CreateEntryResponse](c.httpClient, requestURL, reqJson)
}
