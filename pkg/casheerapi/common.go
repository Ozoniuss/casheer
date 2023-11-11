package casheerapi

import (
	"time"
)

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DefaultLinks allow the user to navigate back to the home page of the API.
type DefaultLinks struct {
	Home string `json:"home"`
}

type HomeLink struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

// MonetaryValueAttributes are used to define the "value" of something in a
// currency agnostic manner.
type MonetaryValueAttributes struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Exponent int    `json:"exponent"`
}

// MonetaryMutableValueAttributes is the same as MonetaryValueAttributes, except
// that it holds the fields that are mutable.
type MonetaryMutableValueAttributes struct {
	Amount   *int    `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Exponent *int    `json:"exponent,omitempty"`
}

// MonetaryValueCreationAttributes is the same as MonetaryValueAttributes,
// but it is used when creating such an object to highlight the optional fields.
type MonetaryValueCreationAttributes struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Exponent *int   `json:"exponent,omitempty"`
}
