package entries

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetEntry(ctx *gin.Context) {

	params, err := common.CtxGetTyped[casheerapi.GetEntryParams](ctx, "queryparams")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	id := ctx.GetInt("entid")

	var entry model.Entry
	err = h.db.WithContext(ctx).Preload("Expenses").Where("id = ?", id).Take(&entry).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.GetEntryResponse{
		Data: EntryToPublic(entry, h.entriesURL, computeRunningTotal(entry.Expenses)),
	}
	includedExpenses := getIncludedExpenses(params.Include, entry.Expenses, h.entriesURL)
	resp.Included = includedExpenses

	ctx.JSON(http.StatusOK, &resp)
}

func getIncludedExpenses(included *string, expenses []model.Expense, entriesURL *url.URL) *[]casheerapi.IncludedExpenseData {
	fmt.Println("heres included", included)
	if included == nil {
		return nil
	}
	includedExpenses := make([]casheerapi.IncludedExpenseData, 0, len(expenses))

	if *included != "expense" {
		return &includedExpenses
	}

	for _, exp := range expenses {
		includedExpenses = append(includedExpenses, IncludedExpensesToPublic(exp, entriesURL))
	}
	return &includedExpenses
}
