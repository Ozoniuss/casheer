package expenses

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/Ozoniuss/casheer/internal/currency"
	ierrors "github.com/Ozoniuss/casheer/internal/errors"
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

func TestHandleListEntry(t *testing.T) {

	entry := newEntry(t, testHandler.db)
	expense1 := model.Expense{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		EntryId:       entry.Id,
		Value:         currency.NewEURValue(1000),
		PaymentMethod: "card",
	}
	expense2 := model.Expense{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		EntryId:       entry.Id,
		Value:         currency.NewEURValue(500),
		PaymentMethod: "cash",
	}
	expense3 := model.Expense{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		EntryId:       entry.Id,
		Value:         currency.NewUSDValue(1000),
		PaymentMethod: "card",
	}
	err := testHandler.db.Create(&[]model.Expense{expense1, expense2, expense3}).Error
	if err != nil {
		t.Fatalf("Could not create expenses: %s\n", err)
	}

	t.Run("Retrieving all expenses of an entry should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("queryparams", casheerapi.ListExpenseParams{})
		ctx.Set("entid", entry.Id)

		testHandler.HandleListExpense(ctx)
		testutils.CheckNoContextErrors(t, ctx)
	})
	t.Run("Retrieving all expenses of a non-existent entry should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("queryparams", casheerapi.ListExpenseParams{})

		ctx.Set("entid", rand.Int())

		testHandler.HandleListExpense(ctx)
		testutils.CheckCanBeContextError(t, ctx, &model.ErrExpenseInvalidEntryKey{})
	})
	t.Run("Retrieving all expenses of an entry should not give an error even when filters are applied", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		currencyFilter := "EUR"
		paymentMethodFilter := "card"
		ctx.Set("queryparams", casheerapi.ListExpenseParams{
			Currency:      &currencyFilter,
			PaymentMethod: &paymentMethodFilter,
		})
		ctx.Set("entid", entry.Id)

		testHandler.HandleListExpense(ctx)
		testutils.CheckNoContextErrors(t, ctx)
	})
}

func TestHandleUpdateEntry(t *testing.T) {

	entry := newEntry(t, testHandler.db)
	expense := newExpense(t, testHandler.db, entry.Id)
	t.Run("Updating an expense with valid fields should update it correctly", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("entid", entry.Id)
		ctx.Set("expid", expense.Id)

		req := casheerapi.UpdateExpenseRequest{
			Amount:        func() *int { m := 600; return &m }(),
			Currency:      func() *string { usd := currency.USD; return &usd }(),
			Exponent:      func() *int { e := 0; return &e }(),
			Name:          func() *string { s := "newname"; return &s }(),
			Description:   func() *string { s := "newdesc"; return &s }(),
			PaymentMethod: func() *string { s := "newpm"; return &s }(),
		}
		ctx.Set("req", req)
		testHandler.HandleUpdateExpense(ctx)

		var savedExpense model.Expense
		err := testHandler.db.Where("id = ?", expense.Id).First(&savedExpense).Error
		if err != nil {
			t.Fatalf("Could not retrieve expense %d: %s\n", savedExpense.Id, err)
		}

		testutils.CheckNoContextErrors(t, ctx)
		if savedExpense.Amount != *req.Amount ||
			savedExpense.Currency != *req.Currency ||
			savedExpense.Exponent != *req.Exponent ||
			savedExpense.Name != *req.Name ||
			savedExpense.Description != *req.Description ||
			savedExpense.PaymentMethod != *req.PaymentMethod {
			t.Errorf("Inserted: %+v\nretrieved %+v\n", req, savedExpense)
		}
	})
	t.Run("Updating an entry with invalid fields should raise an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("entid", entry.Id)
		ctx.Set("expid", expense.Id)

		req := casheerapi.UpdateExpenseRequest{
			Amount:        func() *int { m := 600; return &m }(),
			Currency:      func() *string { usd := "fakecurrency"; return &usd }(),
			Exponent:      func() *int { e := 0; return &e }(),
			Name:          func() *string { s := "newname"; return &s }(),
			Description:   func() *string { s := "newdesc"; return &s }(),
			PaymentMethod: func() *string { s := "newpm"; return &s }(),
		}
		ctx.Set("req", req)
		testHandler.HandleUpdateExpense(ctx)

		var target ierrors.InvalidModel
		testutils.CheckCanBeContextError(t, ctx, &target)
	})

}
