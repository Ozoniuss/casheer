package calls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// makeRequest makes a typed request to the API. It is a helper meant to
// generate the possible requests allowed by the client.
func makeRequest[T ResponseType](method string, client *http.Client, url string, payload []byte) (T, error) {

	var empty T

	var respData T
	var respErr casheerapi.ErrorResponse

	var reqBody io.Reader = nil

	if method == "POST" || method == "PATCH" {
		reqBody = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return empty, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return respData, fmt.Errorf("%s request failed: %w", method, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, fmt.Errorf("reading response body: %w", err)
	}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return empty, fmt.Errorf("unmarshalling response body: %w", err)
	}

	// All those status codes are only issued in case of errors.
	if resp.StatusCode >= 400 {
		err = json.Unmarshal(respBody, &respErr)
		if err != nil {
			return empty, fmt.Errorf("unmarshalling response body: %w", err)
		}
		return empty, respErr
	}

	return respData, nil
}
