package debts

import (
	"net/http"

	"github.com/Ozoniuss/casheer/currency"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleCreateDebt(ctx *gin.Context) {

	req, err := common.CtxGetTyped[casheerapi.CreateDebtRequest](ctx, "req")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	value, err := currency.NewValueBasedOnCurrency(req.Data.Attributes.Amount, req.Data.Attributes.Currency, req.Data.Attributes.Exponent)
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	debt := model.Debt{
		Person:  req.Data.Attributes.Person,
		Value:   value,
		Details: req.Data.Attributes.Details,
	}

	err = h.db.WithContext(ctx).Scopes(model.ValidateModelScope[model.Debt](debt)).Clauses(clause.Returning{}).Create(&debt).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.CreateDebtResponse{
		Data: DebtToPublic(
			debt,
			h.debtsURL,
		),
		Links: common.NewDefaultLinks(h.debtsURL),
	}
	// See https://jsonapi.org/format/#crud-creating-responses-201
	ctx.Header("Locatinon", resp.Data.Links.Self)
	ctx.JSON(http.StatusCreated, &resp)
}
