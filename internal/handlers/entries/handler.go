package entries

import (
	"net/url"

	"gorm.io/gorm"
)

type handler struct {
	db         *gorm.DB
	entriesURL *url.URL
}

func NewHandler(db *gorm.DB, collection *url.URL) handler {
	return handler{
		db:         db,
		entriesURL: collection,
	}
}
