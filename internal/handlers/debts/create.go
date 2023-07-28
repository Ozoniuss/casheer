package debts

import (
	"errors"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	apierrors "github.com/Ozoniuss/casheer/internal/errors"
)

func (h *handler) HandleCreateDebt(ctx *gin.Context) {

	req, ok := common.CtxGetTyped[casheerapi.CreateDebtRequest](ctx, "req")
	if !ok {
		return
	}

	debt := model.Debt{
		Person:  req.Person,
		Amount:  req.Amount,
		Details: req.Details,
	}

	err := h.db.WithContext(ctx).Scopes(model.ValidateModel[model.Debt](debt, model.InvalidDebtErr{})).Clauses(clause.Returning{}).Create(&debt).Error

	// TODO: nicer error handling
	if err != nil {
		switch {
		case errors.Is(err, model.InvalidDebtErr{}):
			{
				common.EmitError(ctx,
					NewInvalidDebtError(err.Error()))
			}
		default:
			common.EmitError(ctx, apierrors.NewUnknownError(err.Error()))
		}
	}

	resp := casheerapi.CreateDebtResponse{
		Data: DebtToPublic(debt, h.apiPaths),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
