package dbstore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ozoniuss/casheer/internal/domain"
	"github.com/Ozoniuss/casheer/internal/domain/currency"
	"github.com/Ozoniuss/casheer/internal/ports/store"
)

type DbStore struct {
	conn *sql.DB
}

func NewDbStore(path string) (*DbStore, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", path, err)
	}

	return &DbStore{
		conn: conn,
	}, nil
}

func (s *DbStore) ListDebts(ctx context.Context) ([]domain.Debt, error) {

	query := `
	SELECT
		id, person, amount, currency, exponent, details, created_at, updated_at
	FROM
		debts
	WHERE
		deleted_at IS NULL
	ORDER BY
		person ASC, amount DESC, id ASC;
	`
	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error running query: %w", err)
	}
	defer rows.Close()
	debts := make([]domain.Debt, 0)

	for rows.Next() {
		var dbdebt DbDebt
		var details sql.NullString

		err := rows.Scan(&dbdebt.ID, &dbdebt.Person, &dbdebt.Amount, &dbdebt.Currency, &dbdebt.Exponent, &details, &dbdebt.CreatedAt, &dbdebt.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning debt row: %w", err)
		}

		dbdebt.Details = details
		debts = append(debts, dbdebt.toDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("looping through rows: %s", err)
	}

	return debts, nil
}
func (s *DbStore) LoadDebt(ctx context.Context, id int) (domain.Debt, error) {

	query := `
	SELECT
		id, person, amount, currency, exponent, details, created_at, updated_at
	FROM
		debts
	WHERE
		id = :id
	AND
		deleted_at IS NULL;
	`
	var dbdebt DbDebt
	var details sql.NullString

	row := s.conn.QueryRowContext(ctx, query, sql.Named("id", id))
	err := row.Scan(&dbdebt.ID, &dbdebt.Person, &dbdebt.Amount, &dbdebt.Currency, &dbdebt.Exponent, &details, &dbdebt.CreatedAt, &dbdebt.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		domainErr := store.ErrNotFound{
			Details: fmt.Sprintf("debt with id %d not found", id),
			Orig:    err,
		}
		return domain.Debt{}, fmt.Errorf("scanning debt row: %w", domainErr)
	case err != nil:
		return domain.Debt{}, fmt.Errorf("scanning debt row: %w", err)
	}

	dbdebt.Details = details
	return dbdebt.toDomain(), nil
}

func (s *DbStore) Healthcheck(ctx context.Context) error {
	return s.conn.PingContext(ctx)
}

func (s *DbStore) Close() error {
	return s.conn.Close()
}

type DbDebt struct {
	ID        int
	Person    string
	Amount    int
	Currency  string
	Exponent  int
	Details   sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d DbDebt) toDomain() domain.Debt {
	return domain.Debt{
		BaseModel: domain.BaseModel{
			Id:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Value: currency.Value{
			Amount:   d.Amount,
			Exponent: d.Exponent,
			Currency: d.Currency,
		},
		Details: d.Details.String,
		Person:  d.Person,
	}
}
