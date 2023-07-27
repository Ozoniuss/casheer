package config

import (
	"fmt"

	cfg "github.com/Ozoniuss/configer"
)

func databaseOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "db-type", Shorthand: "", Value: SQLITE_DB, ConfigKey: "database.type",
			Usage: "Specifies the name of the ports database."},
	}
}

func pgDatabaseOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "pg-db-name", Shorthand: "", Value: "defaultdb", ConfigKey: "postgres_database.name",
			Usage: "Specifies the name of the ports database."},
		{FlagName: "pg-db-host", Shorthand: "", Value: "127.0.0.1", ConfigKey: "postgres_database.host",
			Usage: "Specifies the address on which the ports database listens for connections."},
		{FlagName: "pg-db-port", Shorthand: "", Value: int32(5432), ConfigKey: "postgres_database.port",
			Usage: "Specifies the port on which the ports database listens for connections."},
		{FlagName: "pg-db-user", Shorthand: "", Value: "postgres", ConfigKey: "postgres_database.user",
			Usage: "Specifies the user which connects to the ports database."},
		{FlagName: "pg-db-password", Shorthand: "", Value: "admin", ConfigKey: "postgres_database.password",
			Usage: "Specifies the password of the user which connects to the ports database."},
	}
}

func sqliteDatabaseOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "sqlite-db", Shorthand: "", Value: "casheer.db", ConfigKey: "sqlite_database.file",
			Usage: "Specifies the name of the ports database."},
	}
}

func serverOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "server-address", Shorthand: "", Value: "127.0.0.1", ConfigKey: "server.address",
			Usage: "Specifies the address on which the ports service listens for incoming calls."},
		{FlagName: "server-port", Shorthand: "", Value: int32(8033), ConfigKey: "server.port",
			Usage: "Specifies the port on which the ports service listens for incoming calls."},
	}
}

func apiPathsOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "apipaths-entries", Shorthand: "", Value: "http://localhost:8033/api/entries", ConfigKey: "apipaths.entries",
			Usage: "Specifies the URL of the entry resource collection."},
		{FlagName: "apipaths-totals", Shorthand: "", Value: "http://localhost:8033/api/totals", ConfigKey: "apipaths.totals",
			Usage: "Specifies the URL of the totals resource collection."},
		{FlagName: "apipaths-debts", Shorthand: "", Value: "http://localhost:8033/api/debts", ConfigKey: "apipaths.debts",
			Usage: "Specifies the URL of the debts resource collection."},
	}
}

func allOptions() []cfg.ConfigOption {
	opts := make([]cfg.ConfigOption, 0)
	opts = append(opts, databaseOptions()...)
	opts = append(opts, pgDatabaseOptions()...)
	opts = append(opts, sqliteDatabaseOptions()...)
	opts = append(opts, serverOptions()...)
	opts = append(opts, apiPathsOptions()...)
	return opts
}

func ParseConfig() (Config, error) {
	c := newConfig()

	parserOptions := []cfg.ParserOption{
		cfg.WithConfigName("config"),
		cfg.WithConfigType("yml"),
		cfg.WithConfigPath("./configs"),
		cfg.WithEnvPrefix("CASHEER"),
		cfg.WithEnvKeyReplacer("_"),
		cfg.WithWriteFlag(),
		cfg.WithReadFlag(),
	}

	err := cfg.NewConfig(&c, allOptions(), parserOptions...)
	if err != nil {
		return newConfig(), fmt.Errorf("could not create config: %w", err)
	}
	return c, nil
}
