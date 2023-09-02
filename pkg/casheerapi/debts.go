package casheerapi

const DebtType = "debt"

type DebtData struct {
	ResourceID
	Attributes DebtAtrributes `json:"attributes"`
	Links      DebtLinks      `json:"links"`
}

type DebtListItemData struct {
	ResourceID
	Attributes DebtAtrributes    `json:"attributes"`
	Links      DebtListItemLinks `json:"links"`
}

type DebtAtrributes struct {
	Person  string `json:"person"`
	Amount  int    `json:"amount"`
	Details string `json:"details"`
	Timestamps
}

type DebtLinks struct {
	Self string `json:"self,omitempty"`
}

type DebtListItemLinks struct {
	Self string `json:"self"`
}

type CreateDebtRequest struct {
	// Note that ResourceID and Links are ignored.
	Data DebtData
}

type CreateDebtResponse struct {
	Data  DebtData     `json:"data"`
	Links DefaultLinks `json:"links"`
}

type UpdateDebtData struct {
	Attributes UpdateDebtAttributes `json:"attributes"`
}

type UpdateDebtAttributes struct {
	Person  *string `json:"person,omitempty"`
	Amount  *int    `json:"amount,omitempty"`
	Details *string `json:"details,omitempty"`
}

type UpdateDebtRequest struct {
	Data UpdateDebtData `json:"data"`
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

// Returns an entries link to reveal other possible state transitions.
type ListDebtLinks struct {
	Self    string `json:"self"`
	Entries string `json:"entries"`
}

type ListDebtResponse struct {
	Data  []DebtListItemData `json:"data"`
	Links DefaultLinks       `json:"links"`
}

type GetDebtRequest struct {
}

type GetDebtResponse struct {
	Data DebtData `json:"data"`
}
