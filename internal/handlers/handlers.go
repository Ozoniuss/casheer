package handlers

import "github.com/gin-gonic/gin"

type EntryHandler interface {
	HandleCreateEntry(ctx *gin.Context)
}
