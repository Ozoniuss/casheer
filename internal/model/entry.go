package model

// Entry models an entry of a planning. Fields are automatically mapped by gorm
// to their database equivalents.
type Entry struct {
	BaseModel

	Month         byte
	Year          int16
	Category      string
	Subcategory   string
	ExpectedTotal float32
	RunningTotal  float32
	Recurring     bool
}
