package entries

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/internal/testutils"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mattn/go-sqlite3"
)

var dbname string
var db *gorm.DB
var testHandler handler

const SQL_PATH = "../../../scripts/sqlite/001_tables.up.sql"

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	var err error
	db, dbname, err = testutils.Setup(SQL_PATH)
	if err != nil {
		fmt.Printf("Error setting up tests: %s", err.Error())
		os.Exit(1)
	}
	// Attempt to remove db file anyway.
	defer os.Remove(dbname)
	testHandler = NewHandler(db, &url.URL{
		Scheme: "http",
		Host:   "localhost:69",
		Path:   "/doesnt/matter",
	})
	code := m.Run()
	// call flag.Parse() here if TestMain uses flags
	err = testutils.Teardown(db, dbname)
	if err != nil {
		fmt.Printf("Error cleaning up tests: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(code)
}

func TestHandleCreateEntry(t *testing.T) {

	sharedMonth := 10
	sharedYear := 2023
	sharedEntry := casheerapi.CreateEntryRequest{
		Month:         &sharedMonth,
		Year:          &sharedYear,
		Category:      "category",
		Subcategory:   "subcategory",
		ExpectedTotal: 5000,
	}

	t.Run("Creating a entry should save the entry", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("req", sharedEntry)

		testHandler.HandleCreateEntry(ctx)

		var entrys []model.Entry
		testHandler.db.Find(&entrys)
		if len(entrys) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d", len(entrys))
		}

		savedEntry := entrys[0]
		if savedEntry.Month != *sharedEntry.Month ||
			savedEntry.Year != *sharedEntry.Year ||
			savedEntry.Category != sharedEntry.Category ||
			savedEntry.Subcategory != sharedEntry.Subcategory ||
			savedEntry.ExpectedTotal != sharedEntry.ExpectedTotal {
			t.Errorf("Inserted: %+v\nretrieved %+v", sharedEntry, savedEntry)
		}

	})

	t.Run("Creating an entry with the same month, year, category and subcategory should fail", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("req", sharedEntry)
		testHandler.HandleCreateEntry(ctx)

		var entrys []model.Entry
		testHandler.db.Find(&entrys)
		fmt.Printf("%+v", entrys)
		if len(entrys) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d", len(entrys))
		}
		if len(ctx.Errors) == 0 {
			t.Fatalf("Expected to have an error attached to the context.")
		}
		var ctxerr sqlite3.Error
		ok := errors.As(ctx.Errors[0], &ctxerr)
		if !ok {
			t.Error("Expected error to be of type sqlite3.Error")
		}
		if ctxerr.Code != sqlite3.ErrConstraint && ctxerr.ExtendedCode != sqlite3.ErrConstraintUnique {
			t.Error("Expected error to be Unique Constraint Error")
		}
	})

	t.Run("Creating an invalid entry should yield an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		month := -10 // invalid month
		year := 1900 // invalid year
		dummyEntry := casheerapi.CreateEntryRequest{
			Month:         &month,
			Year:          &year,
			Category:      "category",
			Subcategory:   "subcategory",
			ExpectedTotal: 5000,
		}
		ctx.Set("req", dummyEntry)

		testHandler.HandleCreateEntry(ctx)

		// Should not be in db.
		var entrys []model.Entry
		testHandler.db.Find(&entrys)
		if len(entrys) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d", len(entrys))
		}
		if len(ctx.Errors) == 0 {
			t.Fatalf("Expected to have an error attached to the context.")
		}
		var ctxerr ierrors.InvalidModel
		ok := errors.As(ctx.Errors[0], &ctxerr)
		if !ok {
			t.Error("Expected error to be of type ierrors.InvalidModel")
		}
	})

}

func TestHandleDeleteEntry(t *testing.T) {

	dummyEntry := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         10,
		Year:          2023,
		Category:      "category",
		Subcategory:   "subcategory",
		ExpectedTotal: 5000,
	}

	testHandler.db.Create(&dummyEntry)

	t.Run("Deleting an existing entry should delete the entry", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)
		testHandler.HandleDeleteEntry(ctx)

		var entries []model.Entry
		testHandler.db.Where("id = ?", dummyEntry.Id).Find(&entries)
		if len(entries) != 0 {
			t.Errorf("Expected to have 0 entry, but found %d", len(entries))
		}
	})

	t.Run("Deleting a non-existing entry should fail", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id+1)
		testHandler.HandleDeleteEntry(ctx)

		if len(ctx.Errors) == 0 {
			t.Fatalf("Expected to have an error attached to the context.")
		}
		var ctxerr = gorm.ErrRecordNotFound
		if !errors.Is(ctx.Errors[0], ctxerr) {
			t.Errorf("Expected error to be of type gorm.ErrRecordNotFound, got %s\n", reflect.TypeOf(ctx.Errors[0]))
		}
	})
}
