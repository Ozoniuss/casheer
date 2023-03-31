package expenses

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

func (h *handler) HandleDeleteExpense(ctx *gin.Context) {

	// Todo: maybe index just entry id to not also index expense id?
	entid := ctx.GetString("entid")

	id := ctx.Param("id")
	expuuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewDeleteExpenseFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not delete expense: invalid uuid format: %s", id),
		))
		return
	}

	expense := model.Expense{}
	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(uuid.MustParse(entid))).Clauses(clause.Returning{}).
		Where("id = ?", expuuid).Delete(&expense).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewDeleteExpenseFailedError(
				http.StatusNotFound,
				fmt.Sprintf("Could not delete expense: expense %s not found.", expuuid)))
			return
		default:
			common.EmitError(ctx, NewDeleteExpenseFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not delete expense: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.CreateExpenseResponse{
		Data: ExpenseToPublic(expense),
	}

	ctx.JSON(http.StatusOK, &resp)
}
