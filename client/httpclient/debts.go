package httpclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Ozoniuss/casheer/client/httpclient/calls"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (c *CasheerHTTPClient) CreateDebt(person string, details string, amount int, currency string, exponent int) (public.CreateDebtResponse, error) {
	requestURL := c.debtsURL.String()

	req := public.CreateDebtRequest{
		Data: struct {
			Type       string "json:\"type\" binding:\"required\""
			Attributes struct {
				Person string "json:\"person\" binding:\"required\""
				public.MonetaryValueCreationAttributes
				Details string "json:\"details\""
			} "json:\"attributes\" binding:\"required\""
		}{
			Type: "debt",
			Attributes: struct {
				Person string "json:\"person\" binding:\"required\""
				public.MonetaryValueCreationAttributes
				Details string "json:\"details\""
			}{
				Person:  person,
				Details: details,
				MonetaryValueCreationAttributes: public.MonetaryValueCreationAttributes{
					Currency: currency,
					Amount:   amount,
					Exponent: &exponent,
				},
			},
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return public.CreateDebtResponse{}, fmt.Errorf("marshalling request into JSON: %s", err.Error())
	}

	return calls.MakePOST[public.CreateDebtResponse](c.httpClient, requestURL, reqJson)
}

func (c *CasheerHTTPClient) GetDebt(debtId int) (public.GetDebtResponse, error) {
	requestURL := c.debtsURL.JoinPath(strconv.Itoa(debtId)).String()
	return calls.MakeGET[public.GetDebtResponse](c.httpClient, requestURL)
}

func (c *CasheerHTTPClient) DeleteDebt(debtId int) (public.DeleteDebtResponse, error) {
	requestURL := c.debtsURL.JoinPath(strconv.Itoa(debtId)).String()
	return calls.MakeDELETE[public.DeleteDebtResponse](c.httpClient, requestURL)
}

func (c *CasheerHTTPClient) ListDebts() (public.ListDebtResponse, error) {
	requestURL := c.debtsURL.String()
	return calls.MakeGET[public.ListDebtResponse](c.httpClient, requestURL)
}

func (c *CasheerHTTPClient) ListDebtsForPerson(person string) (public.ListDebtResponse, error) {
	requestURL := &url.URL{}
	*requestURL = *c.debtsURL

	query := requestURL.Query()
	query.Add("person", person)

	requestURL.RawQuery = query.Encode()
	return calls.MakeGET[public.ListDebtResponse](c.httpClient, requestURL.String())
}