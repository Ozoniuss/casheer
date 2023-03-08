package casheerapi

const EntryType = "entry"

type EntryData struct {
	ResourceID
	Month         byte    `json:"month"`
	Year          int16   `json:"year"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	ExpectedTotal float32 `json:"expected_total"`
	RunningTotal  float32 `json:"running_total"`
	Recurring     bool    `json:"recurring"`
	Timestamps
}

type CreateEntryRequest struct {
	Month         *int    `json:"month,omitempty"`
	Year          *int    `json:"year,omitempty"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	ExpectedTotal float32 `json:"expected_total"`
	Recurring     bool    `json:"recurring"`
}

type CreateEntryResponse struct {
	Data EntryData `json:"data"`
}

type UpdateEntryRequest struct {
	Recurring     *bool    `json:"recurring,omitempty"`
	ExpectedTotal *float32 `json:"expected_total,omitempty"`
}

type UpdateEntryResponse struct {
	Data EntryData `json:"data"`
}

type DeleteEntryRequest struct {
}

type DeleteEntryResponse struct {
	Data EntryData `json:"data"`
}
