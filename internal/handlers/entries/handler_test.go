package entries

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"net/url"
	"os"
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
		if savedEntry.Month != *sharedEntry.Month ||
			savedEntry.Year != *sharedEntry.Year ||
			savedEntry.Category != sharedEntry.Category ||
			savedEntry.Subcategory != sharedEntry.Subcategory ||
			savedEntry.ExpectedTotal != sharedEntry.ExpectedTotal {
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
			Month:         &month,
			Year:          &year,
			Category:      "category",
			Subcategory:   "subcategory",
			ExpectedTotal: 5000,
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

	dummyEntry := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         10,
		Year:          2023,
		Category:      "categoryd",
		Subcategory:   "subcategoryd",
		ExpectedTotal: 5000,
	}

	dummyEntryCascade := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         10,
		Year:          2024,
		Category:      "cascade",
		Subcategory:   "cascade",
		ExpectedTotal: 5000,
	}

	dummyExpense := model.Expense{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		EntryId: dummyEntryCascade.Id,
		Value:   100,
		Name:    "dummy expense",
	}

	err := testHandler.db.Create(&[]model.Entry{
		dummyEntry, dummyEntryCascade,
	}).Error
	if err != nil {
		t.Fatalf("Could not create entries: %s\n", err)
	}

	err = testHandler.db.Create(&dummyExpense).Error
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

		ctx.Set("entid", dummyEntry.Id+1)
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

	dummyEntry := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         10,
		Year:          2023,
		Category:      "category2",
		Subcategory:   "subcategory2",
		ExpectedTotal: 5000,
	}
	err := testHandler.db.Create(&dummyEntry).Error
	if err != nil {
		t.Fatalf("Could not create entry: %s\n", err)
	}

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

	dummyEntry1 := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         9,
		Year:          2023,
		Category:      "category3",
		Subcategory:   "subcategory3",
		ExpectedTotal: 3000,
	}
	dummyEntry2 := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         10,
		Year:          2023,
		Category:      "category4",
		Subcategory:   "subcategory4",
		ExpectedTotal: 5000,
	}
	dummyEntry3 := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         11,
		Year:          2023,
		Category:      "category5",
		Subcategory:   "subcategory5",
		ExpectedTotal: 5000,
	}
	err := testHandler.db.Create(&[]model.Entry{dummyEntry1, dummyEntry2, dummyEntry3}).Error
	if err != nil {
		t.Fatalf("Could not create entries: %s\n", err)
	}

	t.Run("Retrieving all entries should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("queryparams", casheerapi.ListEntryParams{})

		testHandler.HandleListEntry(ctx)
		testutils.CheckNoContextErrors(t, ctx)
	})
}

func TestHandleUpdateEntry(t *testing.T) {

	dummyEntry := model.Entry{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Month:         10,
		Year:          2023,
		Category:      "update",
		Subcategory:   "update",
		ExpectedTotal: 5000,
	}
	err := testHandler.db.Create(&dummyEntry).Error
	if err != nil {
		t.Fatalf("Could not create entry: %s\n", err)
	}

	t.Run("Updating an entry with valid fields should update it correctly", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)

		req := casheerapi.UpdateEntryRequest{
			Month:       func() *int { m := 12; return &m }(),
			Year:        func() *int { y := 2024; return &y }(),
			Category:    func() *string { c := "u1"; return &c }(),
			Subcategory: func() *string { s := "u1"; return &s }(),
			Recurring:   func() *bool { r := true; return &r }(),
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
		if savedEntry.Month != *req.Month ||
			savedEntry.Year != *req.Year ||
			savedEntry.Category != *req.Category ||
			savedEntry.Subcategory != *req.Subcategory {
			t.Errorf("Updated: %+v\nretrieved %+v", req, savedEntry)
		}
	})
	t.Run("Updating an entry with invalid fields should raise an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", dummyEntry.Id)

		req := casheerapi.UpdateEntryRequest{
			Month:       func() *int { m := 13; return &m }(),
			Year:        func() *int { y := 2000; return &y }(),
			Category:    func() *string { c := ""; return &c }(),
			Subcategory: func() *string { s := ""; return &s }(),
		}
		ctx.Set("req", req)
		testHandler.HandleUpdateEntry(ctx)

		var target ierrors.InvalidModel
		testutils.CheckCanBeContextError(t, ctx, &target)
	})

}
