package handlers

import "github.com/gin-gonic/gin"

type EntryHandler interface {
	HandleCreateEntry(ctx *gin.Context)
	HandleDeleteEntry(ctx *gin.Context)
	HandleListEntry(ctx *gin.Context)
	HandleUpdateEntry(ctx *gin.Context)
	HandleGetEntry(ctx *gin.Context)
}
