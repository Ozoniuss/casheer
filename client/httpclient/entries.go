package httpclient

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Ozoniuss/casheer/client/httpclient/calls"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (c *CasheerHTTPClient) CreateEntry(month int, year int, category string, subcategory string, expected_total int, currency string, recurring bool) (public.CreateEntryResponse, error) {
	requestURL := c.entriesURL.String()

	req := public.CreateEntryRequest{
		Data: public.CreateEntryData{
			Type: "entry",
			Attributes: public.CreateEntryAttributes{
				Month:       &month,
				Year:        &year,
				Category:    category,
				Subcategory: subcategory,
				Recurring:   recurring,
				ExpectedTotal: public.MonetaryValueCreationAttributes{
					Amount:   expected_total,
					Currency: currency,
					Exponent: nil,
				},
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return public.CreateEntryResponse{}, fmt.Errorf("marshalling request into JSON: %s", err.Error())
	}

	return calls.MakePOST[public.CreateEntryResponse](c.httpClient, requestURL, reqJson)
}

func (c *CasheerHTTPClient) ListEntries() (public.ListEntryResponse, error) {
	requestURL := c.entriesURL.String()
	return calls.MakeGET[public.ListEntryResponse](c.httpClient, requestURL, nil)
}

func (c *CasheerHTTPClient) ListEntriesForPeriod(month, year int) (public.ListEntryResponse, error) {
	requestURL := c.entriesURL.String()
	queryParams := map[string]string{"month": strconv.Itoa(month), "year": strconv.Itoa(year)}
	return calls.MakeGET[public.ListEntryResponse](c.httpClient, requestURL, queryParams)
}

func (c *CasheerHTTPClient) GetEntry(entryId int, includeExpenses bool) (public.GetEntryResponse, error) {
	requestURL := c.entriesURL.JoinPath(strconv.Itoa(entryId)).String()
	queryParams := map[string]string{"include": "expense"}
	if !includeExpenses {
		queryParams = nil
	}
	return calls.MakeGET[public.GetEntryResponse](c.httpClient, requestURL, queryParams)
}
