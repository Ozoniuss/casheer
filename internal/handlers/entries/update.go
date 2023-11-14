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

	var oldEntry model.Entry
	err = h.db.WithContext(ctx).Where("id = ?", id).Take(&oldEntry).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}
	updateEntryFields(req, &oldEntry)

	err = h.db.WithContext(ctx).Preload("Expenses").Clauses(clause.Returning{}).
		Scopes(model.ValidateModel(oldEntry)).Save(&oldEntry).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.UpdateEntryResponse{
		Data: EntryToPublic(oldEntry, h.entriesURL, computeRunningTotal(oldEntry.Expenses)),
	}

	ctx.JSON(http.StatusOK, &resp)
}

func updateEntryFields(req casheerapi.UpdateEntryRequest, entry *model.Entry) {
	if req.Data.Attributes.Month != nil {
		entry.Month = *req.Data.Attributes.Month
	}
	if req.Data.Attributes.Year != nil {
		entry.Year = *req.Data.Attributes.Year
	}
	if req.Data.Attributes.Category != nil {
		entry.Category = *req.Data.Attributes.Category
	}
	if req.Data.Attributes.Subcategory != nil {
		entry.Subcategory = *req.Data.Attributes.Subcategory
	}
	if req.Data.Attributes.Recurring != nil {
		entry.Recurring = *req.Data.Attributes.Recurring
	}
	if req.Data.Attributes.ExpectedTotal.Amount != nil {
		entry.Amount = *req.Data.Attributes.ExpectedTotal.Amount
	}
	if req.Data.Attributes.ExpectedTotal.Exponent != nil {
		entry.Exponent = *req.Data.Attributes.ExpectedTotal.Exponent
	}
	if req.Data.Attributes.ExpectedTotal.Currency != nil {
		entry.Currency = *req.Data.Attributes.ExpectedTotal.Currency
	}
}
