package dbstore

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	"github.com/Ozoniuss/casheer/internal/domain"
	"github.com/Ozoniuss/casheer/internal/domain/currency"
	"golang.org/x/exp/rand"
)

// runMigrations executes the content of the sql file into the database
// instance.
func runMigrations(db *sql.DB, migrations embed.FS) error {
	werr := fs.WalkDir(migrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".down.sql") {
			return nil
		}
		query, err := fs.ReadFile(migrations, path)
		if err != nil {
			return fmt.Errorf("reading sql file: %s", err.Error())
		}
		sqlerr := runSqlQuery(db, string(query))
		if sqlerr != nil {
			return fmt.Errorf("running sql script %s: %s", path, sqlerr.Error())
		}
		return nil
	})
	if werr != nil {
		return fmt.Errorf("going through the sql migration files: %s", werr.Error())
	}
	return nil
}

func runSqlQuery(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("executing sql query: %s", err.Error())
	}
	return nil
}

func newEntry(conn *sql.DB) (domain.Entry, error) {

	month := rand.Intn(12) + 1
	year := rand.Intn(10) + 2023
	category := strconv.Itoa(rand.Int())
	subcategory := strconv.Itoa(rand.Int())

	value := currency.Value{
		Amount:   5000,
		Exponent: -2,
		Currency: "EUR",
	}

	entry, err := domain.NewEntry(month, year, category, subcategory, false, value)
	if err != nil {
		return domain.Entry{}, fmt.Errorf("could not create entry: %s", err.Error())
	}

	_, err = conn.Exec("INSERT INTO entries(month, year, category, subcategory, amount, exponent, currency) VALUES (?,?,?,?,?,?,?);", entry.Month, entry.Year, entry.Category, entry.Subcategory, entry.Amount, entry.Exponent, entry.Currency)
	if err != nil {
		return domain.Entry{}, fmt.Errorf("could not run sql: %s", err.Error())
	}

	return entry, nil
}

func newDebt(conn *sql.DB, name string, id int) (domain.Debt, error) {

	person := name
	value := currency.Value{
		Amount:   5000,
		Exponent: -2,
		Currency: "EUR",
	}

	debt, err := domain.NewDebt(person, "", value)
	if err != nil {
		return domain.Debt{}, fmt.Errorf("could not create debt: %s", err.Error())
	}

	_, err = conn.Exec("INSERT INTO debts(id, person, amount, currency, exponent, details) VALUES (?,?,?,?,?,?);", id, debt.Person, debt.Amount, debt.Currency, debt.Exponent, debt.Details)
	if err != nil {
		return domain.Debt{}, fmt.Errorf("could not run sql: %s", err.Error())
	}

	return debt, nil
}
