package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Ozoniuss/casheer/internal/endpoints"
	"github.com/Ozoniuss/casheer/internal/handlers/debts"
	"github.com/Ozoniuss/casheer/internal/handlers/entries"
	"github.com/Ozoniuss/casheer/internal/handlers/totals"
)

// NewRouter initializes the gin router with the existing handlers and options.
func NewRouter(db *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()
	{
		h := entries.NewHandler(db)
		endpoints.RegisterEntries(r, &h)
	}
	{
		h := debts.NewHandler(db)
		endpoints.RegisterDebts(r, &h)
	}
	{
		h := totals.NewHandler(db)
		endpoints.RegisterTotals(r, &h)
	}
	return r, nil
}
