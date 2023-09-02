package casheerapi

const EntryType = "entry"

// Data returned about an entry when a single one is returned.
type EntryData struct {
	ResourceID
	Attributes    EntryAtrributes    `json:"attributes"`
	Meta          EntryMeta          `json:"meta"`
	Links         EntryLinks         `json:"links"`
	Relationships EntryRelationships `json:"relationships"`
}

type EntryListItemData struct {
	ResourceID
	Attributes    EntryAtrributes    `json:"attributes"`
	Meta          EntryMeta          `json:"meta"`
	Links         EntryListItemLinks `json:"links"`
	Relationships EntryRelationships `json:"relationships"`
}

type EntryMeta struct {
	RunningTotal int `json:"running_total"`
}
type EntryAtrributes struct {
	Month         int    `json:"month"`
	Year          int    `json:"year"`
	Category      string `json:"category"`
	Subcategory   string `json:"subcategory"`
	ExpectedTotal int    `json:"expected_total"`
	Recurring     bool   `json:"recurring"`
	Timestamps
}

type EntryLinks struct {
	Self string `json:"self"`
}

type EntryRelationships struct {
	Expenses EntryExpenseRelationship `json:"expenses"`
}

type EntryExpenseRelationship struct {
	Links EntryExpenseRelationshipLinks `json:"links"`
}

type EntryExpenseRelationshipLinks struct {
	Related string `json:"related"`
}

// No need to return the collection, not the total, as it is returned as a link
// in the listing.
type EntryListItemLinks struct {
	Self string `json:"self"`
}

type CreateEntryRequest struct {
	Month         *int   `json:"month,omitempty"`
	Year          *int   `json:"year,omitempty"`
	Category      string `json:"category"`
	Subcategory   string `json:"subcategory"`
	ExpectedTotal int    `json:"expected_total"`
	Recurring     bool   `json:"recurring"`
}

type CreateEntryResponse struct {
	Data EntryData `json:"data"`
}

type UpdateEntryRequest struct {
	Month         *int    `json:"month,omitempty"`
	Year          *int    `json:"year,omitempty"`
	Category      *string `json:"category,omitempty"`
	Subcategory   *string `json:"subcategory,omitempty"`
	Recurring     *bool   `json:"recurring,omitempty"`
	ExpectedTotal *int    `json:"expected_total,omitempty"`
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

type ListEntryLinks struct {
	Self string   `json:"self"`
	Home HomeLink `json:"home"`
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
