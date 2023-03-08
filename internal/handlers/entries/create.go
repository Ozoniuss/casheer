package entries

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		db: db,
	}
}

func (h *handler) HandleCreateEntry(ctx *gin.Context) {
	ctx.Status(http.StatusCreated)
}
