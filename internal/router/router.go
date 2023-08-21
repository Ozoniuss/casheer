package router

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/endpoints"
	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/Ozoniuss/casheer/internal/handlers/debts"
	"github.com/Ozoniuss/casheer/internal/handlers/entries"
	"github.com/Ozoniuss/casheer/internal/handlers/expenses"
	"github.com/Ozoniuss/casheer/internal/handlers/totals"
	"github.com/Ozoniuss/casheer/internal/middlewares"
)

// NewRouter initializes the gin router with the existing handlers and options.
func NewRouter(db *gorm.DB, config config.Config) (*gin.Engine, error) {
	baseURL := &url.URL{
		Scheme: config.Server.Scheme,
		Host:   fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port),
	}
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())
	{
		h := handlers.NewDefault(config)
		endpoints.RegisterDefaults(r, &h)
	}
	{
		h := entries.NewHandler(db, baseURL.JoinPath(config.ApiPaths.Entries))
		endpoints.RegisterEntries(r, &h)
	}
	{
		h := debts.NewHandler(db, baseURL.JoinPath(config.ApiPaths.Debts))
		endpoints.RegisterDebts(r, &h)
	}
	{
		h := totals.NewHandler(db, config)
		endpoints.RegisterTotals(r, &h)
	}
	{
		h := expenses.NewHandler(db, baseURL.JoinPath(config.ApiPaths.Entries))
		endpoints.RegisterExpenses(r, &h)
	}
	return r, nil
}
