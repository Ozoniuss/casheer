package debts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteDebt(ctx *gin.Context) {

	id := ctx.GetInt("id")

	var debt model.Debt
	err := h.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).Delete(&debt).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewDeleteDebtFailedError(
				http.StatusNotFound,
				fmt.Sprintf("Could not delete debt: Debt %d not found.", id)))
			return
		default:
			common.EmitError(ctx, NewDeleteDebtFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not delete debt: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.CreateDebtResponse{
		Data: DebtToPublic(debt, h.apiPaths),
	}

	resp.Data.Links.Self = ""

	ctx.JSON(http.StatusOK, &resp)
}
