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
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteDebt(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewDeleteDebtFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not delete debt: invalid uuid format: %s", id),
		))
		return
	}

	Debt := model.Debt{}

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", uuid).Delete(&Debt).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewDeleteDebtFailedError(
				http.StatusNotFound,
				fmt.Sprintf("Could not delete Debt: Debt %s not found.", uuid)))
			return
		default:
			common.EmitError(ctx, NewDeleteDebtFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not delete Debt: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.CreateDebtResponse{
		Data: DebtToPublic(Debt),
	}

	ctx.JSON(http.StatusOK, &resp)
}
