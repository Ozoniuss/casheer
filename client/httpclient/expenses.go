package httpclient

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Ozoniuss/casheer/client/httpclient/calls"
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
