package testutils

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Ozoniuss/casheer/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

// Setup creates a temporary sqlite3 database folder, returning a database
// connection.
func Setup(sqlpath string) (*gorm.DB, string, error) {
	var err error
	dbfile, err := os.CreateTemp(".", "*.testdb")
	if err != nil {
		return nil, "", err
	}
	dbname := dbfile.Name()
	dbfile.Close()

	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not open database file: %w", err)
	}

	err = store.RunMigrations(db, sqlpath)
	if err != nil {
		return nil, "", fmt.Errorf("initializing database: %s", err.Error())
	}

	return db, dbname, nil
}

// Teardown closes the sqlite3 database connection, and removes the database
// file.
func Teardown(db *gorm.DB, dbname string) error {
	instance, _ := db.DB()

	// Attempt to close the instance, do nothing in case of error.
	instance.Close()

	err := os.Remove(dbname)
	if err != nil {
		return fmt.Errorf("removing temporary database %s: %s", dbname, err.Error())
	}

	return nil
}
