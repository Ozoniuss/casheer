package model

import "github.com/google/uuid"

// Expense models an expense that can be associated with an entry. Fields are
// automatically mapped by gorm to their database equivalents.
type Expense struct {
	BaseModel

	EntryId       uuid.UUID
	Value         float32 `validate:"required"`
	Name          string  `validate:"required"`
	Description   string
	PaymentMethod string
}
