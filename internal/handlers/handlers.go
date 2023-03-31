package handlers

import "github.com/gin-gonic/gin"

type EntryHandler interface {
	HandleCreateEntry(ctx *gin.Context)
	HandleDeleteEntry(ctx *gin.Context)
	HandleListEntry(ctx *gin.Context)
	HandleUpdateEntry(ctx *gin.Context)
	HandleGetEntry(ctx *gin.Context)
}

type DebtHandler interface {
	HandleCreateDebt(ctx *gin.Context)
	HandleUpdateDebt(ctx *gin.Context)
	HandleDeleteDebt(ctx *gin.Context)
	HandleListDebt(ctx *gin.Context)
	HandleGetDebt(ctx *gin.Context)
}

type ExpenseHandler interface {
	HandleCreateExpense(ctx *gin.Context)
	HandleDeleteExpense(ctx *gin.Context)
	HandleListExpense(ctx *gin.Context)
	HandleUpdateExpense(ctx *gin.Context)
	HandleGetExpense(ctx *gin.Context)
}

type TotalsHandler interface {
	HandleGetRunningTotal(ctx *gin.Context)
}
