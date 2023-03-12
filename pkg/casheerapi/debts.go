package casheerapi

const DebtType = "debt"

type DebtData struct {
	ResourceID
	Person  string  `json:"person"`
	Amount  float32 `json:"amount"`
	Details string  `json:"details"`
	Timestamps
}

type CreateDebtRequest struct {
	Person  string  `json:"person"`
	Amount  float32 `json:"amount"`
	Details string  `json:"details"`
}

type CreateDebtResponse struct {
	Data DebtData `json:"data"`
}

type UpdateDebtRequest struct {
	Person  *string  `json:"person,omitempty"`
	Amount  *float32 `json:"amount,omitempty"`
	Details *string  `json:"details,omitempty"`
}

type UpdateDebtResponse struct {
	Data DebtData `json:"data"`
}

type DeleteDebtRequest struct {
}

type DeleteDebtResponse struct {
	Data DebtData `json:"data"`
}

type ListDebtParams struct {
	Person *string `form:"person,omitempty"`
}

type ListDebtResponse struct {
	Data []DebtData `json:"data"`
}

type GetDebtRequest struct {
}

type GetDebtResponse struct {
	Data DebtData `json:"data"`
}
