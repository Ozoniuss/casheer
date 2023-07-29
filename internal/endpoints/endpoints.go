package endpoints

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/Ozoniuss/casheer/internal/middlewares"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// See following for redirects.
// https://datatracker.ietf.org/doc/html/rfc3986#section-5.4

func RegisterDefaults(router *gin.Engine, h *handlers.DefaultHandler) {
	router.GET("", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/")
	})
	router.GET("/api", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "api/")
	})
	router.GET("/api/", h.DefaultHandler)
}

func RegisterEntries(router *gin.Engine, h handlers.EntryHandler) {
	subrouter := router.Group("/api/entries")

	// Using entid is necessary to avoid conflicting routes with expenses
	// endpoints.
	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateEntryRequest](), h.HandleCreateEntry)
	subrouter.DELETE("/:entid", middlewares.GetURLParam("entid"), h.HandleDeleteEntry)
	subrouter.GET("/", h.HandleListEntry)
	subrouter.PATCH("/:entid", middlewares.GetURLParam("entid"), middlewares.BindJSONRequest[casheerapi.UpdateEntryRequest](), h.HandleUpdateEntry)
	subrouter.GET("/:entid", middlewares.GetURLParam("entid"), h.HandleGetEntry)
}

func RegisterDebts(router *gin.Engine, h handlers.DebtHandler) {
	subrouter := router.Group("/api/debts")

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateDebtRequest](), h.HandleCreateDebt)
	subrouter.DELETE("/:dbtid", middlewares.GetURLParam("dbtid"), h.HandleDeleteDebt)
	subrouter.GET("/", h.HandleListDebt)
	subrouter.PATCH("/:dbtid", middlewares.GetURLParam("dbtid"), middlewares.BindJSONRequest[casheerapi.UpdateDebtRequest](), h.HandleUpdateDebt)
	subrouter.GET("/:dbtid", middlewares.GetURLParam("dbtid"), h.HandleGetDebt)
}

func RegisterExpenses(router *gin.Engine, h handlers.ExpenseHandler) {

	// An expense exists only in the context of an entry. Standalone expenses
	// are not allowed, which the middleware ensures.
	subrouter := router.Group("/api/entries/:entid/expenses/").Use(middlewares.GetURLParam("entid"))

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateExpenseRequest](), h.HandleCreateExpense)
	subrouter.DELETE("/:expid", middlewares.GetURLParam("expid"), h.HandleDeleteExpense)
	subrouter.GET("/", h.HandleListExpense)
	subrouter.PATCH("/:expid", middlewares.GetURLParam("expid"), middlewares.BindJSONRequest[casheerapi.UpdateExpenseRequest](), h.HandleUpdateExpense)
	subrouter.GET("/:expid", middlewares.GetURLParam("expid"), h.HandleGetExpense)
}

func RegisterTotals(router *gin.Engine, h handlers.TotalsHandler) {
	subrouter := router.Group("/api/totals")

	subrouter.GET("/", h.HandleGetRunningTotal)
}
