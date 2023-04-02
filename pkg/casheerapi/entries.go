package casheerapi

const EntryType = "entry"

// Data returned about an entry when a single one is returned.
type EntryData struct {
	ResourceID
	Attributes EntryAtrributes `json:"attributes"`
	Meta       EntryMeta       `json:"meta"`
	Links      EntryLinks      `json:"links"`
}

type EntryListItemData struct {
	ResourceID
	Attributes EntryAtrributes    `json:"attributes"`
	Meta       EntryMeta          `json:"meta"`
	Links      EntryListItemLinks `json:"links"`
}

type EntryMeta struct {
	RunningTotal float32 `json:"running_total"`
}
type EntryAtrributes struct {
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	Category      string  `json:"category"`
	Subcategory   string  `json:"subcategory"`
	ExpectedTotal float32 `json:"expected_total"`
	Recurring     bool    `json:"recurring"`
	Timestamps
}

type EntryLinks struct {
	Self       string `json:"self"`
	Collection string `json:"collection"`
	Expenses   string `json:"expenses"`
	Total      string `json:"total"`
}

// No need to return the collection, not the total, as it is returned as a link
// in the listing.
type EntryListItemLinks struct {
	Self string `json:"self"`
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

// Represents the "total" resource the items listed contribute to. The first
// version of the backend will force a period filter on the listing, meaning
// that only one total can be associated with all elements of an entry
// listing.
//
// Alternatively, a total may be associated with each individual entry, which
// makes sense if entries can be from multiple periods. However, the main use
// case of the application is planning and viewing for a single period (at
// least yet), so this feature doesn't make sense.
type ListEntryLinks struct {
	Self  string `json:"self"`
	Total string `json:"total"`
}

type ListEntryResponse struct {
	Data  []EntryListItemData `json:"data"`
	Links ListEntryLinks      `json:"links"`
}

type GetEntryRequest struct {
}

type GetEntryResponse struct {
	Data EntryData `json:"data"`
}
