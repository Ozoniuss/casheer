package config

//go:generate go tool genconfig -struct=Config -project=Casheer
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
	// CreateIfEmpty is helpful in controlling what the application should do
	// in case an initialized database is not found. In an actual application,
	// this will force using an existing database. However, for testing, this
	// can be enabled in order for the application to create the database and
	// run all migrations.
	//
	// Note that running the migrations requires access to the migration scripts.
	CreateIfEmpty bool
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
