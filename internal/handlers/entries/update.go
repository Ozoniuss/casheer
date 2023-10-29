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
		Scopes(model.ValidateModelScope[model.Entry](oldEntry)).Save(&oldEntry).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
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

func getUpdatedFields(req casheerapi.UpdateEntryRequest) (model.Entry, []string) {

	// See https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make([]string, 0, 6)
	entry := model.Entry{}

	// TODO: proper validation here.
	if req.Data.Attributes.Month != nil {
		updatedFields = append(updatedFields, "month")
		entry.Month = *req.Data.Attributes.Month
	}
	if req.Data.Attributes.Year != nil {
		updatedFields = append(updatedFields, "year")
		entry.Year = *req.Data.Attributes.Year
	}
	if req.Data.Attributes.Category != nil {
		updatedFields = append(updatedFields, "category")
		entry.Category = *req.Data.Attributes.Category
	}
	if req.Data.Attributes.Subcategory != nil {
		updatedFields = append(updatedFields, "subcategory")
		entry.Subcategory = *req.Data.Attributes.Subcategory
	}
	if req.Data.Attributes.Recurring != nil {
		updatedFields = append(updatedFields, "recurring")
		entry.Recurring = *req.Data.Attributes.Recurring
	}
	if req.Data.Attributes.ExpectedTotal.Amount != nil {
		updatedFields = append(updatedFields, "amount")
		entry.Amount = *req.Data.Attributes.ExpectedTotal.Amount
	}
	if req.Data.Attributes.ExpectedTotal.Exponent != nil {
		updatedFields = append(updatedFields, "exponent")
		entry.Exponent = *req.Data.Attributes.ExpectedTotal.Exponent
	}
	if req.Data.Attributes.ExpectedTotal.Currency != nil {
		updatedFields = append(updatedFields, "currency")
		entry.Currency = *req.Data.Attributes.ExpectedTotal.Currency
	}

	return entry, updatedFields
}
