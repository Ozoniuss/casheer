package debts

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteDebt(ctx *gin.Context) {

	id := ctx.GetInt("dbtid")

	var debt model.Debt
	err := h.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).Delete(&debt).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.DeleteDebtResponse{
		Data: DebtToPublic(debt, h.debtsURL),
	}
	resp.Data.Links.Self = ""

	ctx.JSON(http.StatusOK, &resp)
}
