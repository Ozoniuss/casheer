package endpoints

import (
	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/Ozoniuss/casheer/internal/handlers/expenses"
	"github.com/gin-gonic/gin"
)

func RegisterDefaults(router *gin.Engine, h *handlers.DefaultHandler) {
	router.GET("", h.DefaultHandler)
	router.GET("/api", h.DefaultHandler)
	router.GET("/api/", h.DefaultHandler)
}

func RegisterEntries(router *gin.Engine, h handlers.EntryHandler) {
	subrouter := router.Group("/api/entries")

	// Using entid is necessary to avoid conflicting routes with expenses
	// endpoints.
	subrouter.POST("/", h.HandleCreateEntry)
	subrouter.DELETE("/:entid", h.HandleDeleteEntry)
	subrouter.GET("/", h.HandleListEntry)
	subrouter.PATCH("/:entid", h.HandleUpdateEntry)
	subrouter.GET("/:entid", h.HandleGetEntry)
}

func RegisterDebts(router *gin.Engine, h handlers.DebtHandler) {
	subrouter := router.Group("/api/debts")

	subrouter.POST("/", h.HandleCreateDebt)
	subrouter.DELETE("/:id", h.HandleDeleteDebt)
	subrouter.GET("/", h.HandleListDebt)
	subrouter.PATCH("/:id", h.HandleUpdateDebt)
	subrouter.GET("/:id", h.HandleGetDebt)
}

func RegisterExpenses(router *gin.Engine, h handlers.ExpenseHandler) {

	// An expense exists only in the context of an entry. Standalone expenses
	// are not allowed, which the middleware ensures.
	subrouter := router.Group("/api/entries/:entid/expenses/").Use(expenses.RequiredEntryUUID())

	subrouter.POST("/", h.HandleCreateExpense)
	subrouter.DELETE("/:id", h.HandleDeleteExpense)
	subrouter.GET("/", h.HandleListExpense)
	subrouter.PATCH("/:id", h.HandleUpdateExpense)
	subrouter.GET("/:id", h.HandleGetExpense)
}

func RegisterTotals(router *gin.Engine, h handlers.TotalsHandler) {
	subrouter := router.Group("/api/totals")

	subrouter.GET("/", h.HandleGetRunningTotal)
}
