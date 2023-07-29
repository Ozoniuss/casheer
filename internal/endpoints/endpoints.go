package endpoints

import (
	"net/http"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
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

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateEntryRequest](), h.HandleCreateEntry)
	subrouter.POST("", func(ctx *gin.Context) {
		common.EmitError(ctx, ierrors.NewInvalidURLNoTrailingSlashError())
	})

	subrouter.GET("/", h.HandleListEntry)
	subrouter.GET("", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "entries/")
	})
	subrouter.GET("/:entid", middlewares.GetURLParam("entid"), h.HandleGetEntry)

	subrouter.PATCH("/:entid", middlewares.GetURLParam("entid"), middlewares.BindJSONRequest[casheerapi.UpdateEntryRequest](), h.HandleUpdateEntry)
	subrouter.DELETE("/:entid", middlewares.GetURLParam("entid"), h.HandleDeleteEntry)
}

func RegisterDebts(router *gin.Engine, h handlers.DebtHandler) {
	subrouter := router.Group("/api/debts")

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateDebtRequest](), h.HandleCreateDebt)
	subrouter.POST("", func(ctx *gin.Context) {
		common.EmitError(ctx, ierrors.NewInvalidURLNoTrailingSlashError())
	})

	subrouter.GET("/", middlewares.BindQueryParams[casheerapi.ListDebtParams](), h.HandleListDebt)
	subrouter.GET("", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "debts/")
	})
	subrouter.GET("/:dbtid", middlewares.GetURLParam("dbtid"), h.HandleGetDebt)

	subrouter.PATCH("/:dbtid", middlewares.GetURLParam("dbtid"), middlewares.BindJSONRequest[casheerapi.UpdateDebtRequest](), h.HandleUpdateDebt)
	subrouter.DELETE("/:dbtid", middlewares.GetURLParam("dbtid"), h.HandleDeleteDebt)
}

func RegisterExpenses(router *gin.Engine, h handlers.ExpenseHandler) {

	// An expense exists only in the context of an entry. Standalone expenses
	// are not allowed, which the middleware ensures.
	subrouter := router.Group("/api/entries/:entid/expenses").Use(middlewares.GetURLParam("entid"))

	subrouter.POST("/", middlewares.BindJSONRequest[casheerapi.CreateExpenseRequest](), h.HandleCreateExpense)
	subrouter.POST("", func(ctx *gin.Context) {
		common.EmitError(ctx, ierrors.NewInvalidURLNoTrailingSlashError())
	})

	subrouter.GET("/", h.HandleListExpense)
	subrouter.GET("", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "expenses/")
	})
	subrouter.GET("/:expid", middlewares.GetURLParam("expid"), h.HandleGetExpense)

	subrouter.PATCH("/:expid", middlewares.GetURLParam("expid"), middlewares.BindJSONRequest[casheerapi.UpdateExpenseRequest](), h.HandleUpdateExpense)
	subrouter.DELETE("/:expid", middlewares.GetURLParam("expid"), h.HandleDeleteExpense)
}

func RegisterTotals(router *gin.Engine, h handlers.TotalsHandler) {
	subrouter := router.Group("/api/totals")

	subrouter.GET("/", h.HandleGetRunningTotal)
}
