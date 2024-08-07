package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	cfg "github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/router"
	"github.com/Ozoniuss/casheer/internal/store"
)

func run() error {

	log := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{},
	))
	ctx := context.Background()

	config := cfg.NewInitializedConfig()
	log.InfoContext(ctx, "parsed config", slog.Any("config", config))

	conn, err := store.ConnectSqlite(config.SQLiteDatabase.File)
	if err != nil {
		return fmt.Errorf("could not connect to sqlite database: %w", err)
	}

	err = store.EnsureDatabaseFileIsInitialized(conn, config.SQLiteDatabase.File, config.SQLiteDatabase.Migration)
	if err != nil {
		return fmt.Errorf("failed database file check: %w", err)
	}

	log.InfoContext(ctx, "Connected to database.")

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
		fmt.Printf("Error running api: %s", err.Error())
		os.Exit(1)
	}
}
