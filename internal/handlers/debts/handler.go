package debts

import (
	"net/url"

	"github.com/Ozoniuss/casheer/internal/config"
	"gorm.io/gorm"
)

type handler struct {
	db       *gorm.DB
	debtsURL *url.URL
	apiPaths config.ApiPaths
}

func NewHandler(db *gorm.DB, collection *url.URL) handler {
	return handler{
		db:       db,
		debtsURL: collection,
	}
}
