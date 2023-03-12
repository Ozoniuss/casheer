package entries

import (
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) HandleGetEntry(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewDeleteEntryFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not update entry: invalid uuid format: %s", id),
		))
		return
	}

	var entry model.Entry
	err = h.db.Where("id = ?", uuid).Take(&entry).Error

	if err != nil {
		common.EmitError(ctx, NewUpdateEntryFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not retrieve entry: %s", err.Error())))
		return
	}

	resp := casheerapi.GetEntryResponse{
		Data: EntryToPublic(entry),
	}

	ctx.JSON(http.StatusOK, &resp)
}
