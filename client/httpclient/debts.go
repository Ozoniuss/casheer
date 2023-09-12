package httpclient

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Ozoniuss/casheer/client/httpclient/calls"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// CreateDebt(person string, details string, amount string, currency string, exponent string) public.CreateDebtResponse
// DeleteDebt(debtId int) public.DeleteDebtResponse
// UpdateDebt(debtId int, person *string, details *string, amount *string, currency *string, exponent *string) public.UpdateDebtResponse
// GetDebt(debtId int) public.GetDebtResponse
// ListDebts() public.ListDebtResponse
// ListDebtsForPerson(person string) public.ListDebtResponse

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
