package debts

import (
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleCreateDebt(ctx *gin.Context) {

	req, ok := common.CtxGetTyped[casheerapi.CreateDebtRequest](ctx, "req")
	if !ok {
		return
	}

	if req.Person == "" {
		common.EmitError(ctx, NewCreateDebtFailedError(
			http.StatusBadRequest,
			"Cannot create debt: empty person specified."))
		return
	}

	Debt := model.Debt{
		Person:  req.Person,
		Amount:  req.Amount,
		Details: req.Details,
	}

	err := h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&Debt).Error

	// TODO: nicer error handling
	if err != nil {
		common.EmitError(ctx, NewCreateDebtFailedError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not create debt: %s", err.Error())))
		return
	}

	resp := casheerapi.CreateDebtResponse{
		Data: DebtToPublic(Debt, h.apiPaths),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
