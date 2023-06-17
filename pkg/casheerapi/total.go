package casheerapi

const TotalType = "total"

type TotalData struct {
	ResourceID
	Month          int `json:"month"`
	Year           int `json:"year"`
	ExpectedIncome int `json:"expected_income"`
	RunningIncome  int `json:"running_income"`
	ExpectedTotal  int `json:"expected_total"`
	RunningTotal   int `json:"running_total"`
}

type GetTotalParams struct {
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
}

type GetTotalResponse struct {
	Data TotalData `json:"data"`
}
