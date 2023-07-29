package common

import (
	"net/url"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func NewDefaultLinks(baseURL *url.URL) casheerapi.DefaultLinks {

	// Workaround because joining an absolute path actually adds to the
	// end of the path.
	var newUrl url.URL
	newUrl = *baseURL
	newUrl.Path = "/api/"
	return casheerapi.DefaultLinks{
		Home: newUrl.String(),
	}
}
