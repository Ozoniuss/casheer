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

// Error represents an error message that is sent to the client in case the
// HTTP request returns an error.
//
// The error message is modelled after the jsonapi specification:
// https://jsonapi.org/format/#error-objects. Many of the suggested fields had
// not been included, since they were unnecessary for the size of this project.
// In production, however, it is generally a good idea to stick to a
// well-defined standard.
type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
