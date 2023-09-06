package totals

import (
	"github.com/Ozoniuss/casheer/internal/config"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB, config config.Config) handler {
	return handler{
		db: db,
	}
}
