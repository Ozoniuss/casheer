package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Ozoniuss/casheer/internal/endpoints"
	"github.com/Ozoniuss/casheer/internal/handlers/entries"
)

// NewRouter initializes the gin router with the existing handlers and options.
func NewRouter(db *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()

	{
		h := entries.NewHandler(db)
		endpoints.RegisterEntries(r, &h)
	}
	return r, nil
}
