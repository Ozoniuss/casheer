package endpoints

import (
	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/Ozoniuss/casheer/internal/middlewares"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
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
	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateEntryRequest](), h.HandleCreateEntry)
	subrouter.DELETE("/:id", middlewares.GetURLParam("id"), h.HandleDeleteEntry)
	subrouter.GET("/", h.HandleListEntry)
	subrouter.PATCH("/:id", middlewares.GetURLParam("id"), middlewares.BindJSONRequest[casheerapi.UpdateEntryRequest](), h.HandleUpdateEntry)
	subrouter.GET("/:id", middlewares.GetURLParam("id"), h.HandleGetEntry)
}

func RegisterDebts(router *gin.Engine, h handlers.DebtHandler) {
	subrouter := router.Group("/api/debts")

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateDebtRequest](), h.HandleCreateDebt)
	subrouter.DELETE("/:id", middlewares.GetURLParam("id"), h.HandleDeleteDebt)
	subrouter.GET("/", h.HandleListDebt)
	subrouter.PATCH("/:id", middlewares.GetURLParam("id"), middlewares.BindJSONRequest[casheerapi.UpdateDebtRequest](), h.HandleUpdateDebt)
	subrouter.GET("/:id", middlewares.GetURLParam("id"), h.HandleGetDebt)
}

func RegisterExpenses(router *gin.Engine, h handlers.ExpenseHandler) {

	// An expense exists only in the context of an entry. Standalone expenses
	// are not allowed, which the middleware ensures.
	subrouter := router.Group("/api/entries/:entid/expenses/").Use(middlewares.GetURLParam("entid"))

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateExpenseRequest](), h.HandleCreateExpense)
	subrouter.DELETE("/:id", middlewares.GetURLParam("id"), h.HandleDeleteExpense)
	subrouter.GET("/", h.HandleListExpense)
	subrouter.PATCH("/:id", middlewares.GetURLParam("id"), middlewares.BindJSONRequest[casheerapi.UpdateExpenseRequest](), h.HandleUpdateExpense)
	subrouter.GET("/:id", middlewares.GetURLParam("id"), h.HandleGetExpense)
}

func RegisterTotals(router *gin.Engine, h handlers.TotalsHandler) {
	subrouter := router.Group("/api/totals")

	subrouter.GET("/", h.HandleGetRunningTotal)
}
