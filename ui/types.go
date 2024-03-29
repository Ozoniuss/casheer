package main

type DebtListItem struct {
	Id         int
	Person     string
	TotalMoney float32
	Currency   string
	Details    string
}

type CategoryWithEntries struct {
	Category string
	Entries  []EntryListItem
}

type EntryListItem struct {
	Id           int
	TotalMoney   float32
	Currency     string
	Category     string
	Subcategory  string
	Recurring    bool
	RunningTotal map[string]float32
	Expenses     []ExpenseListItem
}

type ExpenseListItem struct {
	Id            int
	TotalMoney    float32
	Currency      string
	Name          string
	Description   string
	PaymentMethod string
}
