package endpoints

import (
	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterEntries(router *gin.Engine, h handlers.EntryHandler) {
	subrouter := router.Group("/api/entries")

	subrouter.POST("/", h.HandleCreateEntry)
	subrouter.DELETE("/:id", h.HandleDeleteEntry)
	subrouter.GET("/", h.HandleListEntry)
	subrouter.PATCH("/:id", h.HandleUpdateEntry)
	subrouter.GET("/:id", h.HandleGetEntry)
}

func RegisterDebts(router *gin.Engine, h handlers.DebtHandler) {
	subrouter := router.Group("/api/debts")

	subrouter.POST("/", h.HandleCreateDebt)
	subrouter.DELETE("/:id", h.HandleDeleteDebt)
	subrouter.GET("/", h.HandleListDebt)
	subrouter.PATCH("/:id", h.HandleUpdateDebt)
	subrouter.GET("/:id", h.HandleGetDebt)
}

func RegisterTotals(router *gin.Engine, h handlers.TotalsHandler) {
	subrouter := router.Group("/api/totals")

	subrouter.GET("/", h.HandleGetRunningTotal)
}
