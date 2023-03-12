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
