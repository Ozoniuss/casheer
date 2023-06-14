package store

import (
	"fmt"
	"time"

	cfg "github.com/Ozoniuss/casheer/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getDSN returns the dsn of the postgres database, as configured by the
// service.
func getDSN(config cfg.PostgresDatabase) string {
	return fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v",
		config.Host, config.Port, config.User, config.Name, config.Password)
}

// ConnectPostgres establishes a connection to the postgres database, using the
// connectivity information provided in the config.
//
// Note that managing the connection pool from go and disabling the database's
// built in pooling system (e.g. pgpool) might give better performance. This is
// not really relevant in my case though.
func ConnectPostgres(config cfg.PostgresDatabase) (*gorm.DB, error) {
	session, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  getDSN(config),
		PreferSimpleProtocol: false, // enables prepared statements (should be active by default)
		WithoutReturning:     false,
	}), &gorm.Config{})

	if err != nil {
		// Note that open doesn't connect to the database, first connection is
		// established when first interacting with it.
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	db, err := session.DB()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve connection pool: %w", err)
	}

	// A session of inserting expenses should not take more than one hour
	// TODO: check how fast the database itself closes an idle connection.
	db.SetConnMaxLifetime(time.Hour)

	// Since I will be the only one using the app I could even lower these
	// to 1...
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	return session, nil
}
