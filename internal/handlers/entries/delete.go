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
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteEntry(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewDeleteEntryFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not delete entry: invalid uuid format: %s", id),
		))
		return
	}

	entry := model.Entry{}

	query := h.db.WithContext(ctx).Clauses(clause.Returning{})
	err = query.Clauses(clause.Returning{}).Where("id = ?", uuid).Delete(&entry).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewDeleteEntryFailedError(
				http.StatusNotFound,
				fmt.Sprintf("Could not delete entry: entry %s not found.", uuid)))
			return
		default:
			common.EmitError(ctx, NewDeleteEntryFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not delete entry: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.CreateEntryResponse{
		Data: EntryToPublic(entry),
	}

	ctx.JSON(http.StatusOK, &resp)
}
