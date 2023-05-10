package expenses

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (h *handler) HandleCreateExpense(ctx *gin.Context) {

	entid := ctx.GetString("entid")

	var req casheerapi.CreateExpenseRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		common.EmitError(ctx, NewCreateExpenseFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	expense := model.Expense{
		EntryId:       uuid.MustParse(entid),
		Value:         req.Value,
		Description:   req.Description,
		Name:          req.Name,
		PaymentMethod: req.PaymentMethod,
	}

	// Todo: custom validator messages
	if err := h.validator.Struct(expense); err != nil {
		common.EmitError(ctx, NewCreateExpenseFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid expense data: %s", err.Error())))
		return
	}

	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(expense.EntryId)).Clauses(clause.Returning{}).Create(&expense).Error

	if err != nil {
		switch {
		case errors.Is(err, &model.NoEntryFoundErr{}):
			common.EmitError(ctx, NewCreateExpenseFailedError(
				http.StatusNotFound,
				fmt.Sprintf("Could not create Expense: no entry with uuid %v", expense.EntryId)))
			return
		default:
			common.EmitError(ctx, NewCreateExpenseFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not create Expense: %s", err.Error())))
			return
		}
	}

	resp := casheerapi.CreateExpenseResponse{
		Data: ExpenseToPublic(expense, h.apiPaths),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
