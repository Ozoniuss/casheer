package store

import (
	"fmt"

	"github.com/Ozoniuss/casheer/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSqlite(config config.SQLiteDatabase) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.File), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open database file: %w", err)
	}

	return db, nil
}
