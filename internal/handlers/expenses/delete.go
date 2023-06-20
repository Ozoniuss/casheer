package expenses

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

func (h *handler) HandleDeleteExpense(ctx *gin.Context) {

	// Todo: maybe index just entry id to not also index expense id?
	entid := ctx.GetInt("entid")
	id := ctx.GetInt("id")

	expense := model.Expense{}
	err := h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).Clauses(clause.Returning{}).
		Where("id = ?", id).Delete(&expense).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewDeleteExpenseFailedError(
				http.StatusNotFound,
				fmt.Sprintf("Could not delete expense: expense %d not found.", id)))
			return
		default:
			common.EmitError(ctx, NewDeleteExpenseFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not delete expense: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.CreateExpenseResponse{
		Data: ExpenseToPublic(expense, h.apiPaths),
	}

	resp.Data.Links.Self = ""
	ctx.JSON(http.StatusOK, &resp)
}
