package casheerapi

const EntryType = "entry"

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
	Month         int                     `json:"month"`
	Year          int                     `json:"year"`
	Category      string                  `json:"category"`
	Subcategory   string                  `json:"subcategory"`
	ExpectedTotal MonetaryValueAttributes `json:"expected_total"`
	Recurring     bool                    `json:"recurring"`
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

type EntryListItemLinks struct {
	Self string `json:"self"`
}

type CreateEntryRequest struct {
	Data CreateEntryData `json:"data"`
}

type CreateEntryData struct {
	Type       string                `json:"type"  `
	Attributes CreateEntryAttributes `json:"attributes"`
}

type CreateEntryAttributes struct {
	Month         *int                            `json:"month,omitempty"`
	Year          *int                            `json:"year,omitempty"`
	Category      string                          `json:"category"  `
	Subcategory   string                          `json:"subcategory"  `
	ExpectedTotal MonetaryValueCreationAttributes `json:"expected_total"  `
	Recurring     bool                            `json:"recurring"`
}

type CreateEntryResponse struct {
	Data EntryData `json:"data"`
}

type UpdateEntryRequest struct {
	Data UpdateEntryData `json:"data"`
}

type UpdateEntryData struct {
	Type       string                `json:"type"`
	Attributes UpdateEntryAttributes `json:"attributes"`
}

type UpdateEntryAttributes struct {
	Month         *int                           `json:"month,omitempty"`
	Year          *int                           `json:"year,omitempty"`
	Category      *string                        `json:"category,omitempty"`
	Subcategory   *string                        `json:"subcategory,omitempty"`
	Recurring     *bool                          `json:"recurring,omitempty"`
	ExpectedTotal MonetaryMutableValueAttributes `json:"expected_total,omitempty"`
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
