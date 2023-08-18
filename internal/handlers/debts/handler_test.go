package debts

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

	_ "github.com/mattn/go-sqlite3"
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

func TestHandleCreateDebt(t *testing.T) {
	t.Run("Creating a debt should save the debt", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		dummyDebt := casheerapi.CreateDebtRequest{
			Person:  "person",
			Amount:  5000,
			Details: "some details",
		}
		ctx.Set("req", dummyDebt)

		testHandler.HandleCreateDebt(ctx)

		testutils.CheckNoContextErrors(t, ctx)

		var debts []model.Debt
		err := testHandler.db.Find(&debts).Error
		if err != nil {
			t.Fatalf("Could not find debts: %s\n", err)
		}
		if len(debts) != 1 {
			t.Errorf("Expected to have 1 debt, but found %d", len(debts))
		}

		savedDebt := debts[0]
		if savedDebt.Amount != dummyDebt.Amount ||
			savedDebt.Details != dummyDebt.Details ||
			savedDebt.Person != dummyDebt.Person {
			t.Errorf("Inserted: %+v\nretrieved %+v", dummyDebt, savedDebt)
		}

	})

	t.Run("Creating the same debt should be ok", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		dummyDebt := casheerapi.CreateDebtRequest{
			Person:  "person",
			Amount:  5000,
			Details: "some details",
		}
		ctx.Set("req", dummyDebt)

		testHandler.HandleCreateDebt(ctx)

		testutils.CheckNoContextErrors(t, ctx)

		var debts []model.Debt
		err := testHandler.db.Find(&debts).Error
		if err != nil {
			t.Fatalf("Could not find debts: %s\n", err)
		}
		if len(debts) != 2 {
			t.Errorf("Expected to have 2 debt, but found %d", len(debts))
		}

		if debts[0].Amount != debts[1].Amount ||
			debts[0].Details != debts[1].Details ||
			debts[0].Person != debts[1].Person {
			t.Errorf("Debts are not the same: %+v and %+v", debts[0], debts[1])
		}
	})

	t.Run("Creating an invalid debt should yield an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		dummyDebt := casheerapi.CreateDebtRequest{
			Person:  "", // invalid person
			Amount:  5000,
			Details: "some details",
		}
		ctx.Set("req", dummyDebt)

		testHandler.HandleCreateDebt(ctx)

		// Should not be in db.
		var debts []model.Debt
		testHandler.db.Find(&debts)
		if len(debts) != 2 {
			t.Errorf("Expected to have 2 debt, but found %d", len(debts))
		}
		testutils.CheckCanBeContextError(t, ctx, &ierrors.InvalidModel{})
	})

}

func TestHandleDeleteDebt(t *testing.T) {
	dummyDebt := model.Debt{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Person:  "person",
		Amount:  5000,
		Details: "some details",
	}
	err := testHandler.db.Create(&dummyDebt).Error
	if err != nil {
		t.Fatalf("Could not create debt:%s\n", err)
	}

	t.Run("Deleting an existing debt should remove it", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("dbtid", dummyDebt.Id)
		testHandler.HandleDeleteDebt(ctx)

		testutils.CheckNoContextErrors(t, ctx)

		var debts []model.Debt
		err := testHandler.db.Where("id = ?", dummyDebt.Id).Find(&debts).Error
		if err != nil {
			t.Fatalf("Could not find debts:%s\n", err)
		}
		if len(debts) != 0 {
			t.Error("Debt did not get deleted.")
		}
	})
}

func TestHandleGetDebt(t *testing.T) {
	dummyDebt := model.Debt{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Person:  "person",
		Amount:  5000,
		Details: "some details",
	}
	err := testHandler.db.Create(&dummyDebt).Error
	if err != nil {
		t.Fatalf("Could not create debt:%s\n", err)
	}

	t.Run("Retrieving an existing debt should not give an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("dbtid", dummyDebt.Id)
		testHandler.HandleGetDebt(ctx)

		testutils.CheckNoContextErrors(t, ctx)
	})

	t.Run("Retrieving an invalid debt should give an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("dbtid", dummyDebt.Id+1)
		testHandler.HandleGetDebt(ctx)

		testutils.CheckIsContextError(t, ctx, gorm.ErrRecordNotFound)
	})
}

func TestHandleListDebt(t *testing.T) {
	dummyDebt := model.Debt{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Person:  "person",
		Amount:  5000,
		Details: "some details",
	}
	err := testHandler.db.Create(&dummyDebt).Error
	if err != nil {
		t.Fatalf("Could not create debt:%s\n", err)
	}

	t.Run("Retrieving all debts should not give an error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("queryparams", casheerapi.ListDebtParams{})
		testHandler.HandleListDebt(ctx)
		testutils.CheckNoContextErrors(t, ctx)
	})
}

func TestHandleUpdateDebt(t *testing.T) {
	dummyDebt := model.Debt{
		BaseModel: model.BaseModel{
			Id: rand.Int(),
		},
		Person:  "person",
		Amount:  5000,
		Details: "some details",
	}
	err := testHandler.db.Create(&dummyDebt).Error
	if err != nil {
		t.Fatalf("Could not update debt: %s\n", err)
	}

	t.Run("Updating an existing debt should update it correctly", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("dbtid", dummyDebt.Id)

		newDebt := casheerapi.UpdateDebtRequest{
			Person:  func() *string { p := "new person"; return &p }(),
			Amount:  func() *int { a := 10000; return &a }(),
			Details: func() *string { d := "new details"; return &d }(),
		}

		ctx.Set("req", newDebt)
		testHandler.HandleUpdateDebt(ctx)

		var debts []model.Debt
		err := testHandler.db.Where("id = ?", dummyDebt.Id).Find(&debts).Error
		if err != nil {
			t.Fatalf("Could not find debts: %s\n", err)
		}

		testutils.CheckNoContextErrors(t, ctx)

		if len(debts) != 1 {
			t.Error("Number of debts is wrong.")
		}

		if debts[0].Amount != *newDebt.Amount ||
			debts[0].Details != *newDebt.Details ||
			debts[0].Person != *newDebt.Person {
			t.Errorf("Debts are not the same: %+v and %+v", debts[0], newDebt)
		}
	})

	t.Run("Updating a debt with invalid fields should raise an error", func(t *testing.T) {

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Set("dbtid", dummyDebt.Id)

		newDebt := casheerapi.UpdateDebtRequest{
			Person:  func() *string { p := ""; return &p }(),
			Amount:  func() *int { a := 10000; return &a }(),
			Details: func() *string { d := "new details"; return &d }(),
		}

		ctx.Set("req", newDebt)
		testHandler.HandleUpdateDebt(ctx)

		var target ierrors.InvalidModel
		testutils.CheckCanBeContextError(t, ctx, &target)

	})
}
