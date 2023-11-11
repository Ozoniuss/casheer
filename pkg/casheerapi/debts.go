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
	Person  string                  `json:"person"`
	Value   MonetaryValueAttributes `json:"value"`
	Details string                  `json:"details"`
	Timestamps
}

type DebtLinks struct {
	Self string `json:"self,omitempty"`
}

type DebtListItemLinks struct {
	Self string `json:"self"`
}

type CreateDebtRequest struct {
	Data CreateDebtData `json:"data"`
}

type CreateDebtData struct {
	Type       string               `json:"type"`
	Attributes CreateDebtAttributes `json:"attributes"`
}

type CreateDebtAttributes struct {
	Value   MonetaryValueCreationAttributes `json:"value"`
	Person  string                          `json:"person"`
	Details string                          `json:"details"`
}

type CreateDebtResponse struct {
	Data  DebtData     `json:"data"`
	Links DefaultLinks `json:"links"`
}

type UpdateDebtRequest struct {
	Data UpdateDebtData `json:"data"  `
}

type UpdateDebtData struct {
	Type       string               `json:"type"`
	Attributes UpdateDebtAttributes `json:"attributes"`
}

type UpdateDebtAttributes struct {
	MonetaryMutableValueAttributes
	Person  *string `json:"person,omitempty"`
	Details *string `json:"details,omitempty"`
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

type ListDebtLinks struct {
	Self string   `json:"self"`
	Home HomeLink `json:"home"`
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
