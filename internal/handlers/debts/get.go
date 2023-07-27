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
)

func (h *handler) HandleGetDebt(ctx *gin.Context) {

	id := ctx.GetInt("dbtid")

	var debt model.Debt
	err := h.db.WithContext(ctx).Where("id = ?", id).Take(&debt).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewGetDebtFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not retrieve debt: debt %d not found.", id)))
			return
		default:
			common.EmitError(ctx, NewGetDebtFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not retrieve debt: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.GetDebtResponse{
		Data: DebtToPublic(debt, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
