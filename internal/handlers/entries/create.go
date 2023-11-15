package entries

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/Ozoniuss/casheer/currency"
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

	value, err := currency.NewValueBasedOnMinorCurrency(req.Data.Attributes.ExpectedTotal.Amount, req.Data.Attributes.ExpectedTotal.Currency, req.Data.Attributes.ExpectedTotal.Exponent)
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	entry := model.Entry{
		Category:    req.Data.Attributes.Category,
		Subcategory: req.Data.Attributes.Subcategory,
		// ExpectedTotal: req.Data.Attributes.ExpectedTotal,
		Value:     model.FromCurrencyValue(value),
		Recurring: req.Data.Attributes.Recurring,
		Month:     int(time.Now().Month()),
		Year:      time.Now().Year(),
	}

	// If month or year are not null, overwrite current month / year.
	if req.Data.Attributes.Month != nil {
		entry.Month = *req.Data.Attributes.Month
	}
	if req.Data.Attributes.Year != nil {
		entry.Year = *req.Data.Attributes.Year
	}

	err = h.db.WithContext(ctx).Scopes(model.ValidateModel(entry)).Clauses(clause.Returning{}).Create(&entry).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.CreateEntryResponse{
		// Running total is obviously 0
		Data: EntryToPublic(entry, h.entriesURL, 0),
	}
	// See https://jsonapi.org/format/#crud-creating-responses-201
	ctx.Header("Locatinon", resp.Data.Links.Self)
	ctx.JSON(http.StatusCreated, &resp)
}
