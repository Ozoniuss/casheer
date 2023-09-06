package handlers

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

type DefaultHandler struct {
}

func NewDefault(config config.Config) DefaultHandler {
	return DefaultHandler{}
}

func (d *DefaultHandler) DefaultHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, casheerapi.PingResponse{
		Info: "Welcome to the casheer api application. Navigate to one of the following links for further action.",
		Links: casheerapi.PingLinks{
			Entries: casheerapi.LinkWithDetails{
				Href:    "entries/",
				Details: "Manage entities.",
			},
			Debts: casheerapi.LinkWithDetails{
				Href:    "debts/",
				Details: "Manage debts.",
			},
			Totals: casheerapi.LinkWithDetails{
				Href:    "totals/",
				Details: "Manage totals.",
			},
		},
	})
}
