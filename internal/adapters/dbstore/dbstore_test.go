package dbstore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Ozoniuss/casheer/internal/ports/store"
	migrations "github.com/Ozoniuss/casheer/scripts"

	_ "github.com/mattn/go-sqlite3"
)

// shared across tests
var testDbStore *DbStore
var testDbFilename string
var exitCode int

func Test_NewDbStore_HealthcheckSuccessful(t *testing.T) {
	err := testDbStore.Healthcheck(context.Background())
	if err != nil {
		t.Errorf("healthcheck should not fail, failed with: %s\n", err.Error())
	}
}

func Test_ListDebts_ReturnsAllDebts_SortedByName_WhenNoFiltersApplied(t *testing.T) {
	defer teardownTest(t, testDbStore)
	createdDebtJohn, err := newDebt(testDbStore.conn, "John", 1)
	if err != nil {
		t.Fatalf("debt was not added: %s\n", err.Error())
	}
	createdDebtAlex, err := newDebt(testDbStore.conn, "Alex", 2)
	if err != nil {
		t.Fatalf("debt was not added: %s\n", err.Error())
	}

	debts, err := testDbStore.ListDebts(context.Background())
	if err != nil {
		t.Fatalf("expected no error when listing debts: %s\n", err.Error())
	}
	if len(debts) != 2 {
		t.Errorf("expected 2 debts, got %d\n", len(debts))
	}

	if createdDebtAlex.Value != debts[0].Value ||
		createdDebtAlex.Details != debts[0].Details ||
		createdDebtAlex.Person != debts[0].Person {
		t.Errorf("debts are different: expected %+v, got %+v\n", createdDebtAlex, debts[0])
	}

	if createdDebtJohn.Value != debts[1].Value ||
		createdDebtJohn.Details != debts[1].Details ||
		createdDebtJohn.Person != debts[1].Person {
		t.Errorf("debts are different: expected %+v, got %+v\n", createdDebtJohn, debts[1])
	}
}

func Test_GetDebt_ReturnsASingleDebt_WhenItExists(t *testing.T) {
	defer teardownTest(t, testDbStore)
	createdDebtJohn, err := newDebt(testDbStore.conn, "John", 1)
	if err != nil {
		t.Fatalf("debt was not added: %s\n", err.Error())
	}
	_, err = newDebt(testDbStore.conn, "Alex", 2)
	if err != nil {
		t.Fatalf("debt was not added: %s\n", err.Error())
	}

	debt, err := testDbStore.LoadDebt(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error when loading debt: %s\n", err.Error())
	}

	if createdDebtJohn.Value != debt.Value ||
		createdDebtJohn.Details != debt.Details ||
		createdDebtJohn.Person != debt.Person {
		t.Errorf("debts are different: expected %+v, got %+v\n", createdDebtJohn, debt)
	}
}

func Test_GetDebt_ReturnsADomainError_WhenTheDebtDoesNotExist(t *testing.T) {
	defer teardownTest(t, testDbStore)
	_, err := newDebt(testDbStore.conn, "John", 1)
	if err != nil {
		t.Fatalf("debt was not added: %s\n", err.Error())
	}

	_, err = testDbStore.LoadDebt(context.Background(), 2)
	if err == nil {
		t.Fatal("expected a non-nil error when loading debt")
	}

	var domainErr store.ErrNotFound
	if !errors.As(err, &domainErr) {
		t.Errorf("expected error to be a domain error, got: %s\n", err.Error())
	}
}

func TestMain(m *testing.M) {
	var err error

	// eagerly defer teardown to avoid resourse leaks
	defer teardown()
	err = setup()
	if err != nil {
		fmt.Printf("error setting up tests database: %s\n", err.Error())
		exitCode = 1
		return
	}
	exitCode = m.Run()
}

func setup() error {
	var err error
	dbfile, err := os.CreateTemp(".", "*.testdb")
	if err != nil {
		return fmt.Errorf("creating temp file: %s", err.Error())
	}
	testDbFilename = dbfile.Name()
	dbfile.Close()

	testDbStore, err = NewDbStore(testDbFilename)
	if err != nil {
		return fmt.Errorf("could not create database connection: %w", err)
	}

	err = runMigrations(testDbStore.conn, migrations.Migrations)
	if err != nil {
		return fmt.Errorf("running migrations: %s", err.Error())
	}

	return nil
}

func teardown() {
	if testDbStore != nil {
		testDbStore.Close()
	}
	if testDbFilename != "" {
		os.Remove(testDbFilename)
	}
	os.Exit(exitCode)
}

func teardownTest(t *testing.T, dbStore *DbStore) {
	_, errDebts := dbStore.conn.Exec("DELETE FROM debts WHERE 1=1;")
	_, errExpenses := dbStore.conn.Exec("DELETE FROM expenses WHERE 1=1;")
	_, errEntries := dbStore.conn.Exec("DELETE FROM entries WHERE 1=1;")

	err := errors.Join(errDebts, errExpenses, errEntries)
	if err != nil {
		t.Fatalf("could not clean up tables: %s\n", err.Error())
	}
}
