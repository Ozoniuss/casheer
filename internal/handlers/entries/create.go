package entries

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (h *handler) HandleCreateEntry(ctx *gin.Context) {

	var req casheerapi.CreateEntryRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, NewCreateEntryFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	if req.Category == "" || req.Subcategory == "" {
		common.EmitError(ctx, NewCreateEntryFailedError(
			http.StatusBadRequest,
			"Category and subcategory cannot be empty."))
		return
	}

	entry := model.Entry{
		Category:      req.Category,
		Subcategory:   req.Subcategory,
		ExpectedTotal: req.ExpectedTotal,
		Recurring:     req.Recurring,
		Month:         int8(time.Now().Month()),
		Year:          int16(time.Now().Year()),
	}

	// If month or year are null, set them to the current month or year.
	if req.Month != nil {
		entry.Month = int8(*req.Month)
	}
	if req.Year != nil {
		entry.Year = int16(*req.Year)
	}

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&entry).Error

	// TODO: nicer error handling
	if err != nil {
		common.EmitError(ctx, NewCreateEntryFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not create entry: %s", err.Error())))
		return
	}

	resp := casheerapi.CreateEntryResponse{
		// Running total is obviously 0
		Data: EntryToPublic(entry, h.apiPaths, 0),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
