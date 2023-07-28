package testutils

import (
	"fmt"
	"os"

	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/internal/store"
	"gorm.io/gorm"
)

// Setup creates a temporary sqlite3 database folder, returning a database
// connection.
func Setup[T model.Entry | model.Debt | model.Expense](table T) (*gorm.DB, string, error) {
	var err error
	dbfile, err := os.CreateTemp(".", "*.db")
	if err != nil {
		return nil, "", err
	}
	dbname := dbfile.Name()
	dbfile.Close()

	db, err := store.ConnectSqlite(dbname)
	if err != nil {
		return nil, "", fmt.Errorf("connecting to temporary database %s: %s", dbname, err.Error())
	}

	// Create the table.
	db.AutoMigrate(table)

	return db, dbname, nil
}

// Teardown closes the sqlite3 database connection, and removes the database
// file.
func Teardown(db *gorm.DB, dbname string) error {
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
