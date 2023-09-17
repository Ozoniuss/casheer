package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Ozoniuss/casheer/client/httpclient/calls"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (c *CasheerHTTPClient) CreateBasicExpense(entryId int, name string, description string, paymentMethod string, amount int, currency string) (public.CreateExpenseResponse, error) {
	requestURL := c.entriesURL.JoinPath(strconv.Itoa(entryId), "expenses/").String()

	req := public.CreateExpenseRequest{
		Data: struct {
			Type       string "json:\"type\" binding:\"required\""
			Attributes struct {
				public.MonetaryValueCreationAttributes
				Name          string "json:\"name\" binding:\"required\""
				Description   string "json:\"description\""
				PaymentMethod string "json:\"payment_method\" binding:\"required\""
			} "json:\"attributes\" binding:\"required\""
		}{
			Type: "expense",
			Attributes: struct {
				public.MonetaryValueCreationAttributes
				Name          string "json:\"name\" binding:\"required\""
				Description   string "json:\"description\""
				PaymentMethod string "json:\"payment_method\" binding:\"required\""
			}{
				MonetaryValueCreationAttributes: public.MonetaryValueCreationAttributes{
					Amount:   amount,
					Currency: currency,
					Exponent: nil,
				},
				Name:          name,
				Description:   description,
				PaymentMethod: paymentMethod,
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return public.CreateExpenseResponse{}, fmt.Errorf("marshalling request into JSON: %s", err.Error())
	}

	return calls.MakePOST[public.CreateExpenseResponse](c.httpClient, requestURL, reqJson)
}

func (c *CasheerHTTPClient) CreateBasicExpenseWithoutId(category string, subcategory string, month int, year int, name string, description string, paymentMethod string, amount int, currency string) (public.CreateExpenseResponse, error) {

	getRequestUrl := &url.URL{}
	*getRequestUrl = *c.entriesURL

	query := getRequestUrl.Query()
	query.Add("month", strconv.Itoa(month))
	query.Add("year", strconv.Itoa(year))
	query.Add("category", category)
	query.Add("subcategory", subcategory)

	getRequestUrl.RawQuery = query.Encode()

	getResp, err := calls.MakeGET[casheerapi.ListEntryResponse](c.httpClient, getRequestUrl.String())
	if err != nil {
		return public.CreateExpenseResponse{}, fmt.Errorf("retrieving entry: %w", err)
	}
	if len(getResp.Data) == 0 {
		return public.CreateExpenseResponse{}, errors.New("entry does not exist")
	}
	if len(getResp.Data) > 1 {
		return public.CreateExpenseResponse{}, errors.New("found multiple entries with the provided data; detected database corruption")
	}

	entryId, err := strconv.Atoi(getResp.Data[0].Id)
	if err != nil {
		return public.CreateExpenseResponse{}, fmt.Errorf("entry id %s is not an int; detected API error", getResp.Data[0].Id)
	}
	return c.CreateBasicExpense(entryId, name, description, paymentMethod, amount, currency)
}
