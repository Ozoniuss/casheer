package store

import (
	"fmt"
	"io"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func createDbFile(dbfile string) error {
	f, err := os.Create(dbfile)
	if err != nil {
		return fmt.Errorf("could not create db file: %s", err.Error())
	}

	return f.Close()
}

// RunMigrations executes the content of the sql file into the database.
func RunMigrations(db *gorm.DB, sqlpath string) error {
	sqlfile, err := os.Open(sqlpath)
	if err != nil {
		return fmt.Errorf("opening sql file: %s", err.Error())
	}
	defer sqlfile.Close()
	query, err := io.ReadAll(sqlfile)
	if err != nil {
		return fmt.Errorf("reading sql file: %s", err.Error())
	}
	err = db.Exec(string(query)).Error
	if err != nil {
		return fmt.Errorf("executing sql query: %s", err.Error())
	}
	return nil
}

func ConnectSqlite(dbfile, sqlpath string) (*gorm.DB, error) {
	shouldRunMigrations := false
	_, err := os.Stat(dbfile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("retrieving stats for %s: %s", dbfile, err)
	} else if os.IsNotExist(err) {
		err := createDbFile(dbfile)
		if err != nil {
			return nil, fmt.Errorf("database file %s doesn't exist and could not create a new one: %s", dbfile, err.Error())
		}
		shouldRunMigrations = true
	}

	db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	if shouldRunMigrations {
		fmt.Println("Initializing database tables...")
		err := RunMigrations(db, sqlpath)
		if err != nil {
			defer os.Remove(dbfile)
			return nil, fmt.Errorf("running sql migrations: %s", err.Error())
		}
	}

	return db, nil
}
