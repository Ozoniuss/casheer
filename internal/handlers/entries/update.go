package entries

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/internal/stringutil"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleUpdateEntry(ctx *gin.Context) {

	id := ctx.GetInt("entid")
	req, ok := common.CtxGetTyped[casheerapi.UpdateEntryRequest](ctx, "req")
	if !ok {
		return
	}

	// Find out what needs to be updated.
	entry, updatedFields := getUpdatedFields(req)

	err := h.db.WithContext(ctx).Preload("Expenses").Select(updatedFields).Clauses(clause.Returning{}).
		Scopes(model.ValidateModelFields[model.Entry](entry, stringutil.CapitalizeArray(updatedFields), model.InvalidEntryErr{})).
		Where("id = ?", id).Updates(&entry).Error

	// TODO: nicer error handling
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewUpdateEntryFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not update entry: entry %d not found.", id)))
		case errors.Is(err, model.InvalidEntryErr{}):
			{
				common.EmitError(ctx, NewCreateEntryFailedError(
					http.StatusBadRequest,
					fmt.Sprintf("Could not update entry: %s", err.Error()),
				))
			}
		default:
			common.EmitError(ctx, NewUpdateEntryFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not update entry: %s", err.Error())))
		}
		return
	}

	resp := casheerapi.UpdateEntryResponse{
		Data: EntryToPublic(entry, h.apiPaths, computeRunningTotal(entry.Expenses)),
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
