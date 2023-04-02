package expenses

import (
	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type handler struct {
	db        *gorm.DB
	validator *validator.Validate
	apiPaths  config.ApiPaths
}

func NewHandler(db *gorm.DB, config config.Config) handler {
	return handler{
		db:        db,
		validator: validator.New(),
		apiPaths:  config.ApiPaths,
	}
}
