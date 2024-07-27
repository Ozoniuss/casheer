package dbstore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ozoniuss/casheer/internal/domain"
	"github.com/Ozoniuss/casheer/internal/domain/currency"
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
	`
	rows, err := s.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error running query: %w", err)
	}
	defer rows.Close()
	debts := make([]domain.Debt, 0)

	for rows.Next() {
		var debt DbDebt
		var details sql.NullString

		err := rows.Scan(&debt.ID, &debt.Person, &debt.Amount, &debt.Currency, &debt.Exponent, &details, &debt.CreatedAt, &debt.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning debt row: %w", err)
		}

		debt.Details = details
		debts = append(debts, debt.toDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("looping through rows: %s", err)
	}

	return debts, nil
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
