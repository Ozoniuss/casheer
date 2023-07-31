package debts

import (
	"net/http"
	"strconv"

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

	debt := model.Debt{
		Person:  req.Person,
		Amount:  req.Amount,
		Details: req.Details,
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
	ctx.Header("Locatinon", strconv.Itoa(debt.Id))
	ctx.JSON(http.StatusCreated, &resp)
}
