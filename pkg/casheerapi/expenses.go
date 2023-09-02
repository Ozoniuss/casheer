package casheerapi

const ExpenseType = "expense"

type ExpenseData struct {
	ResourceID
	Attributes ExpenseAttributes `json:"attributes"`
	Links      ExpenseLinks      `json:"links"`
}

type ExpenseListItemData struct {
	ResourceID
	Attributes ExpenseAttributes    `json:"attributes"`
	Links      ExpenseListItemLinks `json:"links"`
}

type ExpenseAttributes struct {
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Exponent      int    `json:"exponent"`
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
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Exponent      *int   `json:"exponent,omitempty"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PaymentMethod string `json:"payment_method"`
}

type CreateExpenseResponse struct {
	Data ExpenseData `json:"data"`
}

type UpdateExpenseRequest struct {
	Amount        *int    `json:"amount,omitempty"`
	Currency      *string `json:"currency,omitempty"`
	Exponent      *int    `json:"exponent,omitempty"`
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

type ListExpenseLinks struct {
	Self    string `json:"self"`
	Entries string `json:"entries"`
	Debts   string `json:"debts"`
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
