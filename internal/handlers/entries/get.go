package entries

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *handler) HandleGetEntry(ctx *gin.Context) {

	id := ctx.Param("entid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewGetEntryFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not get entry: invalid uuid format: %s", id),
		))
		return
	}

	var entry model.Entry
	err = h.db.WithContext(ctx).Preload("Expenses").Where("id = ?", uuid).Take(&entry).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewGetEntryFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not retrieve entry: entry %s not found.", uuid)))
			return
		default:
			common.EmitError(ctx, NewGetEntryFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not retrieve entry: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.GetEntryResponse{
		Data: EntryToPublic(entry, h.apiPaths, computeRunningTotal(entry.Expenses)),
	}

	ctx.JSON(http.StatusOK, &resp)
}
