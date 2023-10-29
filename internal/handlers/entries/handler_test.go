package entries

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/currency"
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

const SQL_PATH = "../../../scripts/sqlite"

var rand = testutils.NewUniqueRand()

// newEntry creates a random unique entry for this test.
func newEntry(t *testing.T, db *gorm.DB) model.Entry {
	month := rand.Intn(12) + 1
	year := rand.Intn(10) + 2023
	entry := model.Entry{
		Month:       month,
		Year:        year,
		Category:    strconv.Itoa(rand.Int()),
		Subcategory: strconv.Itoa(rand.Int()),
		Value: model.Value{
			Amount:   5000,
			Exponent: -2,
			Currency: "EUR",
		},
	}

	err := testHandler.db.Create(&entry).Error
	if err != nil {
		t.Fatalf("Could not create entry: %s\n", err)
	}
	return entry
}

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
		Data: casheerapi.CreateEntryData{
			Type: "entry",
			Attributes: casheerapi.CreateEntryAttributes{
				Month:       &sharedMonth,
				Year:        &sharedYear,
				Category:    "category",
				Subcategory: "subcategory",
				ExpectedTotal: casheerapi.MonetaryValueCreationAttributes{
					Amount:   5000,
					Currency: "EUR",
				},
			},
		},
	}

	t.Run("Creating a entry with no exponent should save the entry with default exponent", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("req", sharedEntry)

		testHandler.HandleCreateEntry(ctx)

		var entries []model.Entry
		err := testHandler.db.Find(&entries).Error
		if err != nil {
			t.Fatalf("Could not find entry: %s\n", err)
		}

		testutils.CheckNoContextErrors(t, ctx)

		if len(entries) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d\n", len(entries))
		}

		savedEntry := entries[0]
		if savedEntry.Month != *sharedEntry.Data.Attributes.Month ||
			savedEntry.Year != *sharedEntry.Data.Attributes.Year ||
			savedEntry.Category != sharedEntry.Data.Attributes.Category ||
			savedEntry.Subcategory != sharedEntry.Data.Attributes.Subcategory ||
			savedEntry.Amount != sharedEntry.Data.Attributes.ExpectedTotal.Amount ||
			savedEntry.Currency != sharedEntry.Data.Attributes.ExpectedTotal.Currency ||
			savedEntry.Exponent != -2 {
			t.Errorf("Inserted: %+v\nretrieved %+v\n", sharedEntry, savedEntry)
		}

	})

	t.Run("Creating an entry with the same month, year, category and subcategory should fail", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("req", sharedEntry)
		testHandler.HandleCreateEntry(ctx)

		var entries []model.Entry
		err := testHandler.db.Find(&entries).Error
		if err != nil {
			t.Fatalf("Could not find entries: %s\n", err)
		}
		if len(entries) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d", len(entries))
		}

		var target sqlite3.Error
		testutils.CheckCanBeContextError(t, ctx, &target)

		if target.Code != sqlite3.ErrConstraint && target.ExtendedCode != sqlite3.ErrConstraintUnique {
			t.Error("Expected error to be Unique Constraint Error")
		}
	})

	t.Run("Creating an invalid entry should yield an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		month := -10 // invalid month
		year := 1900 // invalid year
		dummyEntry := casheerapi.CreateEntryRequest{
			Data: casheerapi.CreateEntryData{
				Type: "entry",
				Attributes: casheerapi.CreateEntryAttributes{
					Month:       &month,
					Year:        &year,
					Category:    "category",
					Subcategory: "subcategory",
					ExpectedTotal: casheerapi.MonetaryValueCreationAttributes{
						Amount:   5000,
						Currency: "EUR",
					},
				},
			},
		}
		ctx.Set("req", dummyEntry)

		testHandler.HandleCreateEntry(ctx)

		// Should not be in db.
		var entries []model.Entry
		testHandler.db.Find(&entries)
		if len(entries) != 1 {
			t.Errorf("Expected to have 1 entry, but found %d", len(entries))
		}
		var target ierrors.InvalidModel
		testutils.CheckCanBeContextError(t, ctx, &target)
	})

}

func TestHandleDeleteEntry(t *testing.T) {

	dummyEntry := newEntry(t, testHandler.db)
	dummyEntryCascade := newEntry(t, testHandler.db)

	dummyExpense := model.Expense{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		EntryId: dummyEntryCascade.Id,
		Value:   model.FromCurrencyValue(currency.NewUSDValue(1000)),
		Name:    "dummy expense",
	}
	err := testHandler.db.Create(&dummyExpense).Error
	if err != nil {
		t.Fatalf("Could not create expense: %s\n", err)
	}

	t.Run("Deleting an existing entry should delete the entry", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)
		testHandler.HandleDeleteEntry(ctx)

		testutils.CheckNoContextErrors(t, ctx)

		var entries []model.Entry
		testHandler.db.Where("id = ?", dummyEntry.Id).Find(&entries)
		if len(entries) != 0 {
			t.Errorf("Expected to have 0 entry, but found %d", len(entries))
		}
	})

	t.Run("Deleting a non-existing entry should fail", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", rand.Int())
		testHandler.HandleDeleteEntry(ctx)

		testutils.CheckIsContextError(t, ctx, gorm.ErrRecordNotFound)
	})

	t.Run("Deleting an entry with expenses should cascade to expenses", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntryCascade.Id)
		testHandler.HandleDeleteEntry(ctx)

		var expenses []model.Expense
		testHandler.db.Where("entry_id = ?", dummyEntryCascade.Id).Find(&expenses)

		if len(expenses) != 0 {
			t.Errorf("Expected to have 0 expenses, but found %d", len(expenses))
		}
	})
}

func TestHandleGetEntry(t *testing.T) {

	dummyEntry := newEntry(t, testHandler.db)

	t.Run("Retrieving an existing entry should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)
		testHandler.HandleGetEntry(ctx)

		testutils.CheckNoContextErrors(t, ctx)
	})

	t.Run("Retrieving a non-existing entry should give an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id+1)
		testHandler.HandleGetEntry(ctx)

		testutils.CheckIsContextError(t, ctx, gorm.ErrRecordNotFound)
	})
}

func TestHandleListEntry(t *testing.T) {

	newEntry(t, testHandler.db)
	newEntry(t, testHandler.db)
	newEntry(t, testHandler.db)

	t.Run("Retrieving all entries should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("queryparams", casheerapi.ListEntryParams{})

		testHandler.HandleListEntry(ctx)
		testutils.CheckNoContextErrors(t, ctx)
	})
}

func TestHandleUpdateEntry(t *testing.T) {

	dummyEntry := newEntry(t, testHandler.db)

	t.Run("Updating an entry with valid fields should update it correctly", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)

		req := casheerapi.UpdateEntryRequest{
			Data: casheerapi.UpdateEntryData{
				Type: "entry",
				Attributes: casheerapi.UpdateEntryAttributes{
					Month:       func() *int { m := 12; return &m }(),
					Year:        func() *int { y := 2024; return &y }(),
					Category:    func() *string { c := "u1"; return &c }(),
					Subcategory: func() *string { s := "u1"; return &s }(),
					Recurring:   func() *bool { r := true; return &r }(),
				},
			},
		}
		ctx.Set("req", req)
		testHandler.HandleUpdateEntry(ctx)

		var entries []model.Entry
		err := testHandler.db.Where("id = ?", dummyEntry.Id).Find(&entries).Error
		if err != nil {
			t.Fatalf("Could not retrieve entry: %s\n", err)
		}

		testutils.CheckNoContextErrors(t, ctx)

		savedEntry := entries[0]
		if savedEntry.Month != *req.Data.Attributes.Month ||
			savedEntry.Year != *req.Data.Attributes.Year ||
			savedEntry.Category != *req.Data.Attributes.Category ||
			savedEntry.Subcategory != *req.Data.Attributes.Subcategory {
			t.Errorf("Updated: %+v\nretrieved %+v", req, savedEntry)
		}
	})
	t.Run("Updating an entry with invalid fields should raise an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)

		req := casheerapi.UpdateEntryRequest{
			Data: casheerapi.UpdateEntryData{
				Type: "entry",
				Attributes: casheerapi.UpdateEntryAttributes{
					Month:       func() *int { m := 13; return &m }(),
					Year:        func() *int { y := 2000; return &y }(),
					Category:    func() *string { c := ""; return &c }(),
					Subcategory: func() *string { s := ""; return &s }(),
				},
			},
		}
		ctx.Set("req", req)
		testHandler.HandleUpdateEntry(ctx)

		var target ierrors.InvalidModel
		testutils.CheckCanBeContextError(t, ctx, &target)
	})

}
