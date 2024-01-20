package model

type Expense struct {
	Id
	Value

	EntryId       int
	Name          string
	Description   string
	PaymentMethod string
}
