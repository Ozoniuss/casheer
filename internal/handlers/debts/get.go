package debts

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetDebt(ctx *gin.Context) {

	id := ctx.GetInt("dbtid")

	var debt model.Debt
	err := h.db.WithContext(ctx).Where("id = ?", id).Take(&debt).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.GetDebtResponse{
		Data: DebtToPublic(debt, h.debtsURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}
