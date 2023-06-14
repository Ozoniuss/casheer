package config

type Config struct {
	Server           Server
	PostgresDatabase PostgresDatabase `mapstructure:"postgres_database"`
	SQLiteDatabase   SQLiteDatabase   `mapstructure:"sqlite_database"`
	Database         Database
	ApiPaths         ApiPaths
}

const (
	SQLITE_DB   string = "sqlite"
	POSTGRES_DB string = "postgres"
)

type Database struct {
	Type string
}

type PostgresDatabase struct {
	Host     string
	Port     int32
	User     string
	Name     string
	Password string
}

type SQLiteDatabase struct {
	File string
}

type Server struct {
	Address string
	Port    int32
}

type ApiPaths struct {
	Entries string
	Debts   string
	Totals  string
}

func newConfig() Config {
	c := Config{}
	return c
}
