package entries

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetEntry(ctx *gin.Context) {

	id := ctx.GetInt("entid")

	var entry model.Entry
	err := h.db.WithContext(ctx).Preload("Expenses").Where("id = ?", id).Take(&entry).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	includedExpenses := make([]casheerapi.IncludedExpenseData, 0, len(entry.Expenses))
	resp := casheerapi.GetEntryResponse{
		Data: EntryToPublic(entry, h.entriesURL, computeRunningTotal(entry.Expenses)),
	}
	for _, exp := range entry.Expenses {
		includedExpenses = append(includedExpenses, IncludedExpensesToPublic(exp, h.entriesURL))
	}
	resp.Included = &includedExpenses

	ctx.JSON(http.StatusOK, &resp)
}
