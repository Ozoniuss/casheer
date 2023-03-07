package model

// Expense models an expense that can be associated with an entry. Fields are
// automatically mapped by gorm to their database equivalents.
type Expense struct {
	BaseModel

	EntryId       byte
	Value         float32
	Name          string
	Description   string
	PaymentMethod string
}
