package casheerapi

const EntryType = "entry"

type EntryData struct {
	ResourceID
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	ExpectedTotal float32 `json:"expected_total"`
	RunningTotal  float32 `json:"running_total"`
	Recurring     bool    `json:"recurring"`
	Timestamps
	Links EntryLinks `json:"links"`
}

type EntryAtrributes struct {
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	ExpectedTotal float32 `json:"expected_total"`
	RunningTotal  float32 `json:"running_total"`
	Recurring     bool    `json:"recurring"`
	Timestamps
}

type EntryLinks struct {
	Self       string `json:"self"`
	Collection string `json:"collection"`
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
	Month         *int     `json:"month,omitempty"`
	Year          *int     `json:"year,omitempty"`
	Category      *string  `json:"category,omitemtpy"`
	Subcategory   *string  `json:"subcategory,omitempty"`
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

type ListEntryParams struct {
	Month       *int    `form:"month,omitempty"`
	Year        *int    `form:"year,omitempty"`
	Category    *string `form:"category,omitempty"`
	Subcategory *string `form:"subcategory,omitempty"`
}

type ListEntryResponse struct {
	Data []EntryData `json:"data"`
}

type GetEntryRequest struct {
}

type GetEntryResponse struct {
	Data EntryData `json:"data"`
}
