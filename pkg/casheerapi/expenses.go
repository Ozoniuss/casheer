package casheerapi

const ExpenseType = "expense"

type ExpenseData struct {
	ResourceID
	Attributes    ExpenseAttributes    `json:"attributes"`
	Links         ExpenseLinks         `json:"links"`
	Relationships ExpenseRelationships `json:"relationships"`
}

type ExpenseListItemData struct {
	ResourceID
	Attributes ExpenseAttributes    `json:"attributes"`
	Links      ExpenseListItemLinks `json:"links"`
}

type ExpenseAttributes struct {
	MonetaryValueAttributes
	Name          string `json:"name"`
	Description   string `json:"description"`
	PaymentMethod string `json:"payment_method"`
	Timestamps
}

type ExpenseLinks struct {
	Self string `json:"self"`
}

type ExpenseListItemLinks struct {
	Self string `json:"self"`
}

type CreateExpenseRequest struct {
	Data CreateExpenseData `json:"data" binding:"required"`
}

type CreateExpenseData struct {
	Type       string                  `json:"type" binding:"required"`
	Attributes CreateExpenseAttributes `json:"attributes" binding:"required"`
}

type CreateExpenseAttributes struct {
	MonetaryValueCreationAttributes
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type CreateExpenseResponse struct {
	Data ExpenseData `json:"data"`
}

type UpdateExpenseRequest struct {
	Data UpdateExpenseData `json:"data" binding:"required"`
}

type UpdateExpenseData struct {
	Type       string                  `json:"type" binding:"required"`
	Attributes UpdateExpenseAttributes `json:"attributes" binding:"required"`
}
type UpdateExpenseAttributes struct {
	MonetaryMutableValueAttributes
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	PaymentMethod *string `json:"payment_method,omitempty"`
}

type UpdateExpenseResponse struct {
	Data ExpenseData `json:"data"`
}

type DeleteExpenseRequest struct {
}

type DeleteExpenseResponse struct {
	Data ExpenseData `json:"data"`
}

type ListExpenseParams struct {
	AmountGt      *int    `json:"amount[gt],omitempty"`
	AmountLt      *int    `json:"amount[lt],omitempty"`
	Currency      *string `json:"currency,omitempty"`
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	PaymentMethod *string `json:"payment_method,omitempty"`
}

type ListExpenseItemLinks struct {
	Self string `json:"self"`
}

type ListExpenseLinks struct {
	Self string   `json:"self"`
	Home HomeLink `json:"home"`
}

type ExpenseRelationships struct {
	Entries ExpenseEntryRelationship `json:"entries"`
}

type ExpenseEntryRelationship struct {
	Links ExpenseEntryRelationshipLinks `json:"links"`
}

type ExpenseEntryRelationshipLinks struct {
	// Since in this case technically the relationship link is the same as the
	// resource collection link, I've decided to not provide a relationship
	// link at all in order to avoid confusion and be compliant to json:api.
	//
	// The related link simply points to the related expense, and not to a link
	// resource.

	Related string `json:"related"`
}

type ListExpenseResponse struct {
	Data  []ExpenseListItemData `json:"data"`
	Links ListExpenseLinks      `json:"links"`
}

type GetExpenseRequest struct {
}

type GetExpenseResponse struct {
	Data ExpenseData `json:"data"`
}
