package entries

import (
	"net/url"

	"github.com/Ozoniuss/casheer/internal/config"
	"gorm.io/gorm"
)

type handler struct {
	db         *gorm.DB
	entriesURL *url.URL
	apiPaths   config.ApiPaths
}

func NewHandler(db *gorm.DB, collection *url.URL) handler {
	return handler{
		db:         db,
		entriesURL: collection,
	}
}
