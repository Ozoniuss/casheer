package config

type Config struct {
	Server         Server
	SQLiteDatabase SQLiteDatabase `mapstructure:"sqlite_database"`
	ApiPaths       ApiPaths
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

func newConfig() Config {
	c := Config{}
	return c
}
