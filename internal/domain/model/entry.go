package model

type Entry struct {
	Id
	Value

	Month       int
	Year        int
	Category    string
	Subcategory string
	Recurring   bool

	Expenses []Expense
}
