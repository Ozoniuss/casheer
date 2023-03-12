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
	HandleListDebts(ctx *gin.Context)
	HandleGetDebt(ctx *gin.Context)
}
