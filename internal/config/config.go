package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server         Server
	SQLiteDatabase SQLiteDatabase `mapstructure:"sqlite_database"`
	ApiPaths       ApiPaths
}

func NewInitializedConfig() Config {
	return Config{
		Server:         newServerConfig(),
		SQLiteDatabase: newSQLiteDatabaseConfig(),
		ApiPaths:       newAPIPathsConfig(),
	}
}

func newSQLiteDatabaseConfig() SQLiteDatabase {
	return SQLiteDatabase{
		File:      valueOrDefaultString("CASHEER_SQLITE_FILE", "casheer.db"),
		Migration: valueOrDefaultString("CASHEER_SQLITE_MIGRATION", "./scripts/sqlite"),
	}
}

func newServerConfig() Server {
	return Server{
		Scheme:  valueOrDefaultString("CASHEER_SERVER_SCHEME", "http"),
		Address: valueOrDefaultString("CASHEER_SERVER_ADDRESS", "127.0.0.1"),
		Port:    valueOrDefaultInt32("CASHEER_SERVER_PORT", 8033),
	}
}

func newAPIPathsConfig() ApiPaths {
	return ApiPaths{
		Entries:  valueOrDefaultString("CASHEER_APIPATHS_ENTRIES", "entries/"),
		Expenses: valueOrDefaultString("CASHEER_APIPATHS_EXPENSES", "expenses/"),
		Debts:    valueOrDefaultString("CASHEER_APIPATHS_DEBTS", "debts/"),
		Totals:   valueOrDefaultString("CASHEER_APIPATHS_TOTALS", "totals/"),
	}
}

// valueOrDefaultString returns the env var value if it was set, or the
// default value otherwise.
func valueOrDefaultString(envName, def string) string {
	val, set := os.LookupEnv(envName)
	if !set {
		return def
	} else {
		return val
	}
}

// valueOrDefaultInt32 returns the env var value if it was set, or the
// default value otherwise. If the env var value cannot be converted to
// an int32, it returns the default value.
func valueOrDefaultInt32(envName string, def int32) int32 {
	valstr, set := os.LookupEnv(envName)
	if !set {
		return def
	} else {
		val, err := strconv.Atoi(valstr)
		if err != nil {
			return def
		}
		return int32(val)
	}
}

const (
	SQLITE_DB   string = "sqlite"
	POSTGRES_DB string = "postgres"
)

type SQLiteDatabase struct {
	File      string
	Migration string
}

type Server struct {
	Scheme  string
	Address string
	Port    int32
}

type ApiPaths struct {
	Entries  string
	Expenses string
	Debts    string
	Totals   string
}
