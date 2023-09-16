package httpclient

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type CasheerHTTPClient struct {
	httpClient *http.Client
	baseURL    *url.URL
	debtsURL   *url.URL
}

type CasheerClientOpts func(casheerClient *CasheerHTTPClient) error

func WithCustomHTTPClient(customClient *http.Client) CasheerClientOpts {
	return func(casheerClient *CasheerHTTPClient) error {
		if casheerClient.httpClient != nil {
			return errors.New("a custom HTTP client can only be applied once")
		}
		casheerClient.httpClient = customClient
		return nil
	}
}

// WithCustomAuthority allows specifying a custom authority, in the format
// address:port
func WithCustomAuthority(authority string) CasheerClientOpts {
	return func(casheerClient *CasheerHTTPClient) error {
		baseURL, err := url.Parse(fmt.Sprintf("http://%s/api/", authority))
		if err != nil {
			return fmt.Errorf("error creating custom api url with authority %s: %s", authority, err.Error())
		}
		casheerClient.baseURL = baseURL
		return nil
	}
}

func NewCasheerHTTPClient(opts ...CasheerClientOpts) (*CasheerHTTPClient, error) {
	c := &CasheerHTTPClient{}
	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, fmt.Errorf("could not create client: %s", err.Error())
		}
	}
	if c.httpClient == nil {
		c.httpClient = newHTTPClient()
	}
	fmt.Println("base", c.baseURL)
	if c.baseURL == nil {
		c.baseURL, _ = url.Parse("http://localhost:8033/api/")
	}
	c.debtsURL = c.baseURL.JoinPath("debts/")
	return c, nil
}
