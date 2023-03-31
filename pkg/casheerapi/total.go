package casheerapi

const TotalType = "total"

type TotalData struct {
	ResourceID
	Month          int     `json:"month"`
	Year           int     `json:"year"`
	ExpectedIncome float32 `json:"expected_income"`
	RunningIncome  float32 `json:"running_income"`
	ExpectedTotal  float32 `json:"expected_total"`
	RunningTotal   float32 `json:"running_total"`
}

type GetTotalParams struct {
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
}

type GetTotalResponse struct {
	Data TotalData `json:"data"`
}
