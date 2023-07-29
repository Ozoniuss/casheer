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
	Self       string `json:"self,omitempty"`
	Collection string `json:"collection"`
}

type DebtListItemLinks struct {
	Self string `json:"self"`
}

type CreateDebtRequest struct {
	Person  string `json:"person"`
	Amount  int    `json:"amount"`
	Details string `json:"details"`
}

type CreateDebtResponse struct {
	Data  DebtData     `json:"data"`
	Links DefaultLinks `json:"links"`
}

type UpdateDebtRequest struct {
	Person  *string `json:"person,omitempty"`
	Amount  *int    `json:"amount,omitempty"`
	Details *string `json:"details,omitempty"`
}

type UpdateDebtResponse struct {
	Data DebtData `json:"data"`
}

type DeleteDebtRequest struct {
}

type DeleteDebtResponse struct {
	Data  DebtData  `json:"data"`
	Links DebtLinks `json:"links"`
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
	Links ListDebtLinks      `json:"links"`
}

type GetDebtRequest struct {
}

type GetDebtResponse struct {
	Data DebtData `json:"data"`
}
