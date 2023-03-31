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
)

func (h *handler) HandleGetExpense(ctx *gin.Context) {

	// Todo: maybe index just entry id to not also index expense id?
	entid := ctx.GetString("entid")

	// Todo: middleware for parsing uuid parameters.
	id := ctx.Param("id")
	expuuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewGetExpenseFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not get expense: invalid uuid format: %s", id),
		))
		return
	}

	var expense model.Expense
	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(uuid.MustParse(entid))).
		Where("id = ?", expuuid).Take(&expense).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewGetExpenseFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not retrieve expense: expense %s not found.", expuuid)))
			return
		default:
			common.EmitError(ctx, NewGetExpenseFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not retrieve expense: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.GetExpenseResponse{
		Data: ExpenseToPublic(expense),
	}

	ctx.JSON(http.StatusOK, &resp)
}
