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

func (h *handler) HandleUpdateExpense(ctx *gin.Context) {

	entid := ctx.GetString("entid")

	id := ctx.Param("id")
	expuuid, err := uuid.Parse(id)
	if err != nil {
		common.EmitError(ctx, NewUpdateExpenseFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not update Expense: invalid uuid format: %s", id),
		))
		return
	}

	var req casheerapi.UpdateExpenseRequest
	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, NewUpdateExpenseFailed(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
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
	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(uuid.MustParse(entid))).
		Model(&expense).Clauses(clause.Returning{}).Where("id = ?", expuuid).Updates(updatedFields).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, NewUpdateExpenseFailed(
				http.StatusNotFound,
				fmt.Sprintf("Could not update expense: expense %s not found.", expuuid)))
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
