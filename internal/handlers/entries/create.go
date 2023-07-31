package entries

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (h *handler) HandleCreateEntry(ctx *gin.Context) {

	req, err := common.CtxGetTyped[casheerapi.CreateEntryRequest](ctx, "req")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	entry := model.Entry{
		Category:      req.Category,
		Subcategory:   req.Subcategory,
		ExpectedTotal: req.ExpectedTotal,
		Recurring:     req.Recurring,
		Month:         int(time.Now().Month()),
		Year:          time.Now().Year(),
	}

	// If month or year are not null, overwrite current month / year.
	if req.Month != nil {
		entry.Month = *req.Month
	}
	if req.Year != nil {
		entry.Year = *req.Year
	}

	err = h.db.WithContext(ctx).Scopes(model.ValidateModelScope[model.Entry](entry)).Clauses(clause.Returning{}).Create(&entry).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.CreateEntryResponse{
		// Running total is obviously 0
		Data: EntryToPublic(entry, h.entriesURL, 0),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
