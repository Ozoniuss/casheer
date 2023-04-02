package main

import (
	"fmt"
	"os"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/router"
	"github.com/Ozoniuss/casheer/internal/store"

	log "github.com/Ozoniuss/stdlog"
)

func run() error {
	config, err := config.ParseConfig()
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}
	log.Infof("parsed config: %+v\n", config)

	conn, err := store.Connect(config.Database)
	if err != nil {
		return fmt.Errorf("could not connect to postgres database: %w", err)
	}
	log.Infoln("Connected to database.")

	engine, err := router.NewRouter(conn, config)
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
