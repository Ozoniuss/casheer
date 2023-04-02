package handlers

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

type DefaultHandler struct {
	apiPaths config.ApiPaths
}

func NewDefault(config config.Config) DefaultHandler {
	return DefaultHandler{
		apiPaths: config.ApiPaths,
	}
}

func (d *DefaultHandler) DefaultHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, casheerapi.PingResponse{
		Info: "Welcome to the casheer api application. Navigate to one of the following links for further action.",
		Links: casheerapi.PingLinks{
			Entries: casheerapi.LinkWithDetails{
				Href:    d.apiPaths.Entries + "/",
				Details: "Manage entities.",
			},
			Debts: casheerapi.LinkWithDetails{
				Href:    d.apiPaths.Debts + "/",
				Details: "Manage debts.",
			},
		},
	})
}
