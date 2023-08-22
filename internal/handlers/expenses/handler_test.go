package expenses

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/internal/currency"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/internal/testutils"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dbname string
var db *gorm.DB
var testHandler handler

const SQL_PATH = "../../../scripts/sqlite/001_tables.up.sql"

var rand = testutils.NewUniqueRand()

// newEntry creates a random unique entry for this test.
func newEntry(t *testing.T, db *gorm.DB) model.Entry {
	month := rand.Intn(12) + 1
	year := rand.Intn(10) + 2023
	entry := model.Entry{
		Month:         month,
		Year:          year,
		Category:      strconv.Itoa(rand.Int()),
		Subcategory:   strconv.Itoa(rand.Int()),
		ExpectedTotal: 5000,
	}

	err := testHandler.db.Create(&entry).Error
	if err != nil {
		t.Fatalf("Could not create entry: %s\n", err)
	}
	return entry
}

// newExpense creates a random expense for this test.
func newExpense(t *testing.T, db *gorm.DB, entryId int) model.Expense {
	expense := model.Expense{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		EntryId:       entryId,
		Value:         currency.NewEURValue(100),
		Name:          "myexpense",
		Description:   "mydescription",
		PaymentMethod: "card",
	}

	err := testHandler.db.Create(&expense).Error
	if err != nil {
		t.Fatalf("Could not create expense: %s\n", err)
	}
	return expense
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

func TestHandleCreateExpense(t *testing.T) {

	entry := newEntry(t, db)

	t.Run("Creating a valid expense should save the expense", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)

		req := casheerapi.CreateExpenseRequest{
			Amount:        100,
			Currency:      "EUR",
			Name:          "test",
			Description:   "test",
			PaymentMethod: "card",
		}
		ctx.Set("req", req)

		testHandler.HandleCreateExpense(ctx)

		var expenses []model.Expense
		err := testHandler.db.Find(&expenses).Error
		if err != nil {
			t.Fatalf("Could not find expense: %s\n", err)
		}

		testutils.CheckNoContextErrors(t, ctx)

		if len(expenses) != 1 {
			t.Errorf("Expected to have 1 expense, but found %d\n", len(expenses))
		}

		savedExpense := expenses[0]
		if savedExpense.Amount != req.Amount ||
			savedExpense.Currency != req.Currency ||
			savedExpense.Exponent != -2 ||
			savedExpense.Name != req.Name ||
			savedExpense.Description != req.Description ||
			savedExpense.PaymentMethod != req.PaymentMethod {
			t.Errorf("Inserted: %+v\nretrieved %+v\n", req, savedExpense)
		}

	})

	t.Run("Creating an expense with invalid currency should fail", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)

		req := casheerapi.CreateExpenseRequest{
			Amount:        100,
			Currency:      "invalid",
			Name:          "test",
			Description:   "test",
			PaymentMethod: "card",
		}
		ctx.Set("req", req)

		testHandler.HandleCreateExpense(ctx)

		var target currency.ErrInvalidCurrency
		testutils.CheckCanBeContextError(t, ctx, &target)
	})

	t.Run("Creating an expense with invalid entryId should fail", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", rand.Int())

		req := casheerapi.CreateExpenseRequest{
			Amount:        100,
			Currency:      "EUR",
			Name:          "test",
			Description:   "test",
			PaymentMethod: "card",
		}
		ctx.Set("req", req)

		testHandler.HandleCreateExpense(ctx)

		var target model.ErrExpenseInvalidEntryKey
		testutils.CheckCanBeContextError(t, ctx, &target)
	})

}

func TestHandleDeleteEntry(t *testing.T) {

	entry := newEntry(t, testHandler.db)
	expense := newExpense(t, testHandler.db, entry.Id)

	t.Run("Deleting an existing expense should delete the expense", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)
		ctx.Set("expid", expense.Id)

		testHandler.HandleDeleteExpense(ctx)
		testutils.CheckNoContextErrors(t, ctx)

		var expenses []model.Expense
		testHandler.db.Where("id = ?", expense.Id).Find(&expenses)
		if len(expenses) != 0 {
			t.Errorf("Expected to have 0 entry, but found %d", len(expenses))
		}
	})

	t.Run("Deleting an expense with invalid ID should fail", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)
		ctx.Set("expid", rand.Int())

		testHandler.HandleDeleteExpense(ctx)
		testutils.CheckIsContextError(t, ctx, gorm.ErrRecordNotFound)
	})

	t.Run("Deleting an existing expense with invalid entry ID should fail", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", rand.Int())
		ctx.Set("expid", expense.Id)

		testHandler.HandleDeleteExpense(ctx)
		testutils.CheckCanBeContextError(t, ctx, &model.ErrExpenseInvalidEntryKey{})
	})

}

func TestHandleGetEntry(t *testing.T) {

	entry := newEntry(t, testHandler.db)
	expense := newExpense(t, testHandler.db, entry.Id)

	t.Run("Retrieving an existing expense should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)
		ctx.Set("expid", expense.Id)
		testHandler.HandleGetExpense(ctx)

		testutils.CheckNoContextErrors(t, ctx)
	})

	t.Run("Retrieving a non-existing entry with valid entry id should give an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)
		ctx.Set("expid", rand.Int())
		testHandler.HandleGetExpense(ctx)

		testutils.CheckIsContextError(t, ctx, gorm.ErrRecordNotFound)
	})

	t.Run("Retrieving any expense with non-existing entry id should give an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", rand.Int())
		ctx.Set("expid", expense.Id)
		testHandler.HandleGetExpense(ctx)

		testutils.CheckCanBeContextError(t, ctx, &model.ErrExpenseInvalidEntryKey{})
	})
}

// func TestHandleListEntry(t *testing.T) {

// 	newEntry(t, testHandler.db)
// 	newEntry(t, testHandler.db)
// 	newEntry(t, testHandler.db)

// 	t.Run("Retrieving all entries should not give an error", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)
// 		ctx.Set("queryparams", casheerapi.ListEntryParams{})

// 		testHandler.HandleListEntry(ctx)
// 		testutils.CheckNoContextErrors(t, ctx)
// 	})
// }

// func TestHandleUpdateEntry(t *testing.T) {

// 	dummyEntry := newEntry(t, testHandler.db)

// 	t.Run("Updating an entry with valid fields should update it correctly", func(t *testing.T) {

// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)

// 		ctx.Set("entid", dummyEntry.Id)

// 		req := casheerapi.UpdateEntryRequest{
// 			Month:       func() *int { m := 12; return &m }(),
// 			Year:        func() *int { y := 2024; return &y }(),
// 			Category:    func() *string { c := "u1"; return &c }(),
// 			Subcategory: func() *string { s := "u1"; return &s }(),
// 			Recurring:   func() *bool { r := true; return &r }(),
// 		}
// 		ctx.Set("req", req)
// 		testHandler.HandleUpdateEntry(ctx)

// 		var entries []model.Entry
// 		err := testHandler.db.Where("id = ?", dummyEntry.Id).Find(&entries).Error
// 		if err != nil {
// 			t.Fatalf("Could not retrieve entry: %s\n", err)
// 		}

// 		testutils.CheckNoContextErrors(t, ctx)

// 		savedEntry := entries[0]
// 		if savedEntry.Month != *req.Month ||
// 			savedEntry.Year != *req.Year ||
// 			savedEntry.Category != *req.Category ||
// 			savedEntry.Subcategory != *req.Subcategory {
// 			t.Errorf("Updated: %+v\nretrieved %+v", req, savedEntry)
// 		}
// 	})
// 	t.Run("Updating an entry with invalid fields should raise an error", func(t *testing.T) {

// 		w := httptest.NewRecorder()
// 		ctx, _ := gin.CreateTestContext(w)

// 		ctx.Set("entid", dummyEntry.Id)

// 		req := casheerapi.UpdateEntryRequest{
// 			Month:       func() *int { m := 13; return &m }(),
// 			Year:        func() *int { y := 2000; return &y }(),
// 			Category:    func() *string { c := ""; return &c }(),
// 			Subcategory: func() *string { s := ""; return &s }(),
// 		}
// 		ctx.Set("req", req)
// 		testHandler.HandleUpdateEntry(ctx)

// 		var target ierrors.InvalidModel
// 		testutils.CheckCanBeContextError(t, ctx, &target)
// 	})

// }
