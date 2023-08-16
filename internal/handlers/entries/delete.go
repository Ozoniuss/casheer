package entries

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteEntry(ctx *gin.Context) {

	id := ctx.GetInt("entid")

	var entry model.Entry
	var expenses []model.Expense

	// workaround because cascading soft deletes doesn't fucking work
	err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txerr := h.db.Where("entry_id = ?", id).Find(&expenses).Error
		if txerr != nil {
			return txerr
		}

		txerr = h.db.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&entry).Error
		if txerr != nil {
			return txerr
		}
		// Preload doesn't fucking work either, but having the expenses is
		// required in order to produce the output resource.
		entry.Expenses = expenses

		txerr = h.db.Delete(&expenses).Error
		if txerr != nil {
			return txerr
		}

		return nil
	})
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
