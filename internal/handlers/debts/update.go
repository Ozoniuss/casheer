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

func (h *handler) HandleUpdateDebt(ctx *gin.Context) {

	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewUpdateDebtFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not update debt: invalid uuid format: %s", id),
		))
		return
	}

	var req casheerapi.UpdateDebtRequest
	err = ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, NewUpdateDebtFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make(map[string]any)

	if req.Person != nil {
		updatedFields["person"] = req.Person
	}
	if req.Amount != nil {
		updatedFields["amount"] = *req.Amount
	}
	if req.Details != nil {
		updatedFields["details"] = req.Details
	}

	var Debt model.Debt
	err = h.db.WithContext(ctx).Model(&Debt).Clauses(clause.Returning{}).Where("id = ?", uuid).Updates(updatedFields).Error

	// TODO: nicer error handling
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewUpdateDebtFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not update debt: debt %s not found.", uuid)))
			return
		default:
			common.EmitError(ctx, NewUpdateDebtFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not update debt: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.UpdateDebtResponse{
		Data: DebtToPublic(Debt, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
