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
	Person string `json:"person"`
	MonetaryValueAttributes
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
	Data struct {
		Type       string `json:"type" binding:"required"`
		Attributes struct {
			Person string `json:"person" binding:"required"`
			MonetaryValueCreationAttributes
			Details string `json:"details"`
		} `json:"attributes" binding:"required"`
	} `json:"data" binding:"required"`
}

type CreateDebtResponse struct {
	Data  DebtData     `json:"data"`
	Links DefaultLinks `json:"links"`
}

type UpdateDebtRequest struct {
	Data struct {
		Type       string `json:"type" binding:"required"`
		Attributes struct {
			Person  *string `json:"person,omitempty"`
			Details *string `json:"details,omitempty"`
			MonetaryMutableValueAttributes
		} `json:"attributes" binding:"required"`
	} `json:"data" binding:"required"`
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
