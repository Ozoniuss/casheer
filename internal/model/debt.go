package model

// Debt models a debt owed to or held by someone. Fields are automatically
// mapped by gorm to their database equivalents.
type Debt struct {
	BaseModel

	Person  string
	Amount  float32
	Details string
}
