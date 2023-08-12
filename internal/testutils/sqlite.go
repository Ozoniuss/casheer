package testutils

import (
	"fmt"
	"io"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	sqlfile, err := os.Open(sqlpath)
	if err != nil {
		return nil, "", fmt.Errorf("opening sql file: %s", err.Error())
	}
	query, err := io.ReadAll(sqlfile)
	if err != nil {
		return nil, "", fmt.Errorf("reading sql file: %s", err.Error())
	}
	err = db.Exec(string(query)).Error
	if err != nil {
		return nil, "", fmt.Errorf("executing sql query: %s", err.Error())
	}

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
