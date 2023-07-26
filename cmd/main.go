package main

import (
	"fmt"
	"os"

	cfg "github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/router"
	"github.com/Ozoniuss/casheer/internal/store"
	"gorm.io/gorm"

	log "github.com/Ozoniuss/stdlog"
)

func run() error {
	config, err := cfg.ParseConfig()
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}
	log.Infof("parsed config: %+v\n", config)

	var conn *gorm.DB

	switch config.Database.Type {
	case cfg.SQLITE_DB:
		conn, err = store.ConnectSqlite(config.SQLiteDatabase.File)
		if err != nil {
			return fmt.Errorf("could not connect to sqlite database: %w", err)
		}
	case cfg.POSTGRES_DB:
		conn, err = store.ConnectPostgres(config.PostgresDatabase)
		if err != nil {
			return fmt.Errorf("could not connect to postgres database: %w", err)
		}
	default:
		return fmt.Errorf("invalid database specified: %s", config.Database.Type)
	}

	log.Infoln("Connected to database.")

	engine, err := router.NewRouter(conn, config)
	if err != nil {
		return fmt.Errorf("could not create new router: %w", err)
	}
	engine.Run(fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port))

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Errf("Error running api: %s", err.Error())
		os.Exit(1)
	}
}
