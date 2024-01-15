package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ozoniuss/casheer/client/httpclient"
	"github.com/Ozoniuss/casheer/internal/store"
	"gorm.io/gorm"
)

var casheerClient, _ = httpclient.NewCasheerHTTPClient(
	httpclient.WithCustomAuthority("localhost:6597"),
)

var conn *gorm.DB

func TestMain(m *testing.M) {
	dbname := os.Getenv("DBNAME")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not get working directory: %s", err.Error())
	}

	conn, err = store.ConnectSqlite(filepath.Join(wd, "..", dbname))
	if err != nil {
		fmt.Printf("could not connect to sqlite database: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}
