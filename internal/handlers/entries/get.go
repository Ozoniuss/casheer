package entries

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) HandleGetEntry(ctx *gin.Context) {

	id := ctx.GetInt("entid")

	var entry model.Entry
	err := h.db.WithContext(ctx).Preload("Expenses").Where("id = ?", id).Take(&entry).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewGetEntryFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not retrieve entry: entry %d not found.", id)))
			return
		default:
			common.EmitError(ctx, NewGetEntryFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not retrieve entry: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.GetEntryResponse{
		Data: EntryToPublic(entry, h.entriesURL, computeRunningTotal(entry.Expenses)),
	}

	ctx.JSON(http.StatusOK, &resp)
}
