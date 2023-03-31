package expenses

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type handler struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		db:        db,
		validator: validator.New(),
	}
}
