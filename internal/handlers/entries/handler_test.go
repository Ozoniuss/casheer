package entries

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/internal/store"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

var dbname string
var db *gorm.DB
var testHandler handler

func setup() error {
	var err error
	dbfile, err := os.CreateTemp(".", "*.db")
	if err != nil {
		return err
	}
	dbname = dbfile.Name()
	dbfile.Close()

	db, err = store.ConnectSqlite(dbname)
	if err != nil {
		return fmt.Errorf("connecting to temporary database %s: %s", dbname, err.Error())
	}

	// Create the entries table.
	db.AutoMigrate(model.Entry{})

	testHandler = NewHandler(db, config.Config{})
	gin.SetMode(gin.TestMode)
	return nil
}

func teardown() error {
	instance, _ := db.DB()

	err := instance.Close()
	if err != nil {
		return fmt.Errorf("closing testing database %s: %s", dbname, err.Error())
	}

	err = os.Remove(dbname)
	if err != nil {
		return fmt.Errorf("removing temporary database %s: %s", dbname, err.Error())
	}
	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Printf("Error setting up tests: %s", err.Error())
		os.Exit(1)
	}
	// Attempt to remove file anyway.
	defer os.Remove(dbname)
	code := m.Run()
	// call flag.Parse() here if TestMain uses flags
	err = teardown()
	if err != nil {
		fmt.Printf("Error cleaning up tests: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(code)
}

func TestHandleCreateExpense(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	month := 10
	year := 2021
	ctx.Set("req", casheerapi.CreateEntryRequest{
		Month:         &month,
		Year:          &year,
		Category:      "category",
		Subcategory:   "subcategory",
		ExpectedTotal: 1000,
	})

	testHandler.HandleCreateEntry(ctx)

	t.Run("Response format", func(t *testing.T) {
		var resp casheerapi.CreateEntryResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			panic(err)
		}
		// check attributes.
	})

	t.Run("Number of elements in database", func(t *testing.T) {
		var entries []model.Entry
		db.Find(&entries)
		if len(entries) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d", len(entries))
		}
	})
}
