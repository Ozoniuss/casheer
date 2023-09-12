package e2e

import (
	"fmt"
	"os"
	"testing"

	"github.com/Ozoniuss/casheer/client/httpclient"
)

var casheerClient *httpclient.CasheerHTTPClient

func TestMain(m *testing.M) {

	var err error
	casheerClient, err = httpclient.NewCasheerHTTPClient(
		httpclient.WithCustomAuthority("localhost:8033"),
	)
	if err != nil {
		fmt.Printf("Could not create casheer client: %s\n", err.Error())
	}
	os.Exit(m.Run())
}
