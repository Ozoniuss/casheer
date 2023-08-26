package entries

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleUpdateEntry(ctx *gin.Context) {

	id := ctx.GetInt("entid")
	req, err := common.CtxGetTyped[casheerapi.UpdateEntryRequest](ctx, "req")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	entry, updatedFields := getUpdatedFields(req)
	err = h.db.WithContext(ctx).Preload("Expenses").Select(updatedFields).Clauses(clause.Returning{}).
		Scopes(model.ValidateModelScope[model.Entry](entry)).
		Where("id = ?", id).Updates(&entry).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
	}

	resp := casheerapi.UpdateEntryResponse{
		Data: EntryToPublic(entry, h.entriesURL, computeRunningTotal(entry.Expenses)),
	}

	ctx.JSON(http.StatusOK, &resp)
}

func getUpdatedFields(req casheerapi.UpdateEntryRequest) (model.Entry, []string) {

	// See https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make([]string, 0, 6)
	entry := model.Entry{}

	// TODO: proper validation here.
	if req.Month != nil {
		updatedFields = append(updatedFields, "month")
		entry.Month = *req.Month
	}
	if req.Year != nil {
		updatedFields = append(updatedFields, "year")
		entry.Year = *req.Year
	}
	if req.Category != nil {
		updatedFields = append(updatedFields, "category")
		entry.Category = *req.Category
	}
	if req.Subcategory != nil {
		updatedFields = append(updatedFields, "subcategory")
		entry.Subcategory = *req.Subcategory
	}
	if req.Recurring != nil {
		updatedFields = append(updatedFields, "recurring")
		entry.Recurring = *req.Recurring
	}
	if req.ExpectedTotal != nil {
		updatedFields = append(updatedFields, "expected_total")
		entry.ExpectedTotal = *req.ExpectedTotal
	}

	return entry, updatedFields
}
