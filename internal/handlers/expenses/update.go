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

func (h *handler) HandleUpdateExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	id := ctx.GetInt("id")

	req, ok := common.CtxGetTyped[casheerapi.UpdateExpenseRequest](ctx, "req")
	if !ok {
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make(map[string]any)

	if req.Value != nil {
		updatedFields["value"] = *req.Value
	}
	if req.Description != nil {
		updatedFields["description"] = *req.Description
	}
	if req.Name != nil {
		updatedFields["name"] = *req.Name
	}
	if req.PaymentMethod != nil {
		updatedFields["payment_method"] = *req.PaymentMethod
	}

	var expense model.Expense
	err := h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).
		Model(&expense).Clauses(clause.Returning{}).Where("id = ?", id).Updates(updatedFields).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewUpdateExpenseFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not update expense: expense %d not found.", id)))
			return
		default:
			common.EmitError(ctx, NewUpdateExpenseFailed(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not update expense: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.UpdateExpenseResponse{
		Data: ExpenseToPublic(expense, h.apiPaths),
	}

	ctx.JSON(http.StatusOK, &resp)
}
