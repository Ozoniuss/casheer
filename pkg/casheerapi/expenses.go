package casheerapi

const ExpenseType = "expense"

type ExpenseData struct {
	ResourceID
	Value         float32 `json:"value"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	PaymentMethod string  `json:"payment_method"`
	Timestamps
}

type CreateExpenseRequest struct {
	Value         float32 `json:"value"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	PaymentMethod string  `json:"payment_method"`
}

type CreateExpenseResponse struct {
	Data ExpenseData `json:"data"`
}

type UpdateExpenseRequest struct {
	Value         *float32 `json:"value,omitempty"`
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	PaymentMethod *string  `json:"payment_method,omitempty"`
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
	Value         *float32 `json:"value,omitempty"`
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	PaymentMethod *string  `json:"payment_method,omitempty"`
}

type ListExpenseResponse struct {
	Data []ExpenseData `json:"data"`
}

type GetExpenseRequest struct {
}

type GetExpenseResponse struct {
	Data ExpenseData `json:"data"`
}
