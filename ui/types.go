package main

type DebtListItem struct {
	Id         int
	Person     string
	TotalMoney float32
	Currency   string
	Details    string
}

type EntryListItem struct {
	Id           int
	TotalMoney   float32
	Currency     string
	Category     string
	Subcategory  string
	Recurring    bool
	RunningTotal float32
}
