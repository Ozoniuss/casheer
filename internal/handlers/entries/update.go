package entries

import (
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleUpdateEntry(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewUpdateEntryFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not update entry: invalid uuid format: %s", id),
		))
		return
	}

	var req casheerapi.UpdateEntryRequest
	err = ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, NewUpdateEntryFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make(map[string]any)

	if req.Month != nil {
		updatedFields["month"] = int8(*req.Month)
	}
	if req.Year != nil {
		updatedFields["year"] = int16(*req.Year)
	}
	if req.Category != nil {
		updatedFields["category"] = *req.Category
	}
	if req.Subcategory != nil {
		updatedFields["subcategory"] = *req.Subcategory
	}
	if req.Recurring != nil {
		updatedFields["recurring"] = *req.Recurring
	}
	if req.ExpectedTotal != nil {
		updatedFields["expected_total"] = *req.ExpectedTotal
	}

	var entry model.Entry
	err = h.db.WithContext(ctx).Model(&entry).Clauses(clause.Returning{}).Where("id = ?", uuid).Updates(updatedFields).Error

	// TODO: nicer error handling
	if err != nil {
		common.EmitError(ctx, NewUpdateEntryFailed(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not update entry: %s", err.Error())))
		return
	}

	resp := casheerapi.UpdateEntryResponse{
		Data: EntryToPublic(entry),
	}

	ctx.JSON(http.StatusOK, &resp)
}
