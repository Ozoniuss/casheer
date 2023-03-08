package endpoints

import (
	"github.com/Ozoniuss/casheer/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterEntries(router *gin.Engine, h handlers.EntryHandler) {
	subrouter := router.Group("/api/entries")

	subrouter.GET("/:id", h.HandleCreateEntry)
}
