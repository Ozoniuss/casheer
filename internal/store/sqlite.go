package store

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSqlite(dbfile string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	return db, nil
}
