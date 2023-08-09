package entries

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteEntry(ctx *gin.Context) {

	id := ctx.GetInt("entid")

	var entry model.Entry
	err := h.db.WithContext(ctx).Clauses(clause.Returning{}).Preload("expenses").Where("id = ?", id).Delete(&entry).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.CreateEntryResponse{
		Data: EntryToPublic(entry, h.entriesURL, computeRunningTotal(entry.Expenses)),
	}
	resp.Data.Links.Self = ""

	ctx.JSON(http.StatusOK, &resp)
}
