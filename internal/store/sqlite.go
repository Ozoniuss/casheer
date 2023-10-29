package store

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

// runSql runs the sql script from the file at the given path.
func runSql(tx *gorm.DB, path string) error {
	sqlfile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening sql file: %s", err.Error())
	}
	defer sqlfile.Close()
	query, err := io.ReadAll(sqlfile)
	if err != nil {
		return fmt.Errorf("reading sql file: %s", err.Error())
	}
	err = tx.Exec(string(query)).Error
	if err != nil {
		return fmt.Errorf("executing sql query: %s", err.Error())
	}
	return nil
}

// RunMigrations executes the content of the sql file into the database.
func RunMigrations(db *gorm.DB, migrationDir string) error {

	dir, err := os.Stat(migrationDir)
	if err != nil {
		return fmt.Errorf("reading migrations directory %s: %s", migrationDir, err.Error())
	}
	if !dir.IsDir() {
		return fmt.Errorf("%s it not a directory", migrationDir)
	}

	werr := filepath.WalkDir(migrationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".down.sql") {
			return nil
		}
		sqlerr := runSql(db, path)
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

func ConnectSqlite(dbfile, sqlpath string) (*gorm.DB, error) {
	shouldRunMigrations := false
	fstat, err := os.Stat(dbfile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("retrieving stats for %s: %s", dbfile, err)
	} else if os.IsNotExist(err) {
		err := createDbFile(dbfile)
		if err != nil {
			return nil, fmt.Errorf("database file %s doesn't exist and could not create a new one: %s", dbfile, err.Error())
		}
		shouldRunMigrations = true
	} else if err == nil && fstat.Size() == 0 {
		// file exists, but has size 0. Should still run the migrations.
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
