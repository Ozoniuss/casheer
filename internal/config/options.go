package config

import (
	"fmt"

	cfg "github.com/Ozoniuss/configer"
)

func databaseOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "db-name", Shorthand: "", Value: "defaultdb", ConfigKey: "database.name",
			Usage: "Specifies the name of the ports database"},
		{FlagName: "db-host", Shorthand: "", Value: "127.0.0.1", ConfigKey: "database.host",
			Usage: "Specifies the address on which the ports database listens for connections"},
		{FlagName: "db-port", Shorthand: "", Value: int32(5432), ConfigKey: "database.port",
			Usage: "Specifies the port on which the ports database listens for connections"},
		{FlagName: "db-user", Shorthand: "", Value: "postgres", ConfigKey: "database.user",
			Usage: "Specifies the user which connects to the ports database"},
		{FlagName: "db-password", Shorthand: "", Value: "admin", ConfigKey: "database.password",
			Usage: "Specifies the password of the user which connects to the ports database"},
	}
}

func serverOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "server-address", Shorthand: "", Value: "127.0.0.1", ConfigKey: "server.address",
			Usage: "Specifies the address on which the ports service listens for incoming calls"},
		{FlagName: "server-port", Shorthand: "", Value: int32(8033), ConfigKey: "server.port",
			Usage: "Specifies the port on which the ports service listens for incoming calls"},
	}
}

func allOptions() []cfg.ConfigOption {
	opts := make([]cfg.ConfigOption, 0)
	opts = append(opts, databaseOptions()...)
	opts = append(opts, serverOptions()...)
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
	}

	err := cfg.NewConfig(&c, allOptions(), parserOptions...)
	if err != nil {
		return newConfig(), fmt.Errorf("could not create config: %w", err)
	}
	return c, nil
}
