package entries

import (
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandleListEntry(ctx *gin.Context) {

	var params casheerapi.ListEntryParams
	err := ctx.BindQuery(&params)

	if err != nil {
		common.EmitError(ctx, NewListEntryFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	if params.Month != nil {
		filters["month"] = int8(*params.Month)
	}
	if params.Year != nil {
		filters["year"] = int16(*params.Year)
	}
	if params.Category != nil {
		filters["category"] = *params.Category
	}
	if params.Subcategory != nil {
		filters["subcategory"] = *params.Subcategory
	}

	var entries []model.Entry
	err = h.db.WithContext(ctx).Where(filters).Order("year desc").Order("month desc").Find(&entries).Error

	// TODO: nicer error handling
	if err != nil {
		common.EmitError(ctx, NewListEntryFailedError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list entries: %s", err.Error())))
		return
	}

	resp := casheerapi.ListEntryResponse{
		Data: make([]casheerapi.EntryData, 0, len(entries)),
	}
	for _, e := range entries {
		resp.Data = append(resp.Data, EntryToPublic(e))
	}
	ctx.JSON(http.StatusOK, &resp)
}
