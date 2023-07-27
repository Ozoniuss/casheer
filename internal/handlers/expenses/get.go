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
)

func (h *handler) HandleGetExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	id := ctx.GetInt("expid")

	var expense model.Expense
	err := h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).
		Where("id = ?", id).Take(&expense).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewGetExpenseFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not retrieve expense: expense %d not found.", id)))
			return
		default:
			common.EmitError(ctx, NewGetExpenseFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not retrieve expense: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.GetExpenseResponse{
		Data: ExpenseToPublic(expense, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
