package casheerapi

import "time"

type ResourceID struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
