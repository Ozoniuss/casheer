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
