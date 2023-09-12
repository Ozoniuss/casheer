package httpclient

import (
	"net/http"
	"time"
)

// newHTTPClient initializes a new HTTP client, used under the hood by the
// casheer client to make requests to the server. It is meant to be used as a
// singleton because of the client's ability to cache TCP connection and reuse
// them accross multiple requests.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
