package casheerapi

const TotalType = "total"

type TotalData struct {
	ResourceID
	Month          int   `json:"month"`
	Year           int   `json:"year"`
	ExpectedIncome int64 `json:"expected_income"`
	RunningIncome  int64 `json:"running_income"`
	ExpectedTotal  int64 `json:"expected_total"`
	RunningTotal   int64 `json:"running_total"`
}

type GetTotalParams struct {
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
}

type GetTotalResponse struct {
	Data TotalData `json:"data"`
}
