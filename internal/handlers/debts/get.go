package debts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (h *handler) HandleGetDebt(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewGetDebtFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not update debt: invalid uuid format: %s", id),
		))
		return
	}

	var debt model.Debt
	err = h.db.WithContext(ctx).Where("id = ?", uuid).Take(&debt).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewGetDebtFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not retrieve debt: debt %s not found.", uuid)))
			return
		default:
			common.EmitError(ctx, NewGetDebtFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not retrieve debt: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.GetDebtResponse{
		Data: DebtToPublic(debt),
	}

	ctx.JSON(http.StatusOK, &resp)
}
