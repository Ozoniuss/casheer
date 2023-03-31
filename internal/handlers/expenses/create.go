package expenses

import (
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

	expid := ctx.GetString("entid")

	var req casheerapi.CreateExpenseRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		common.EmitError(ctx, NewCreateExpenseFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind request body: %s", err.Error())))
		return
	}

	expense := model.Expense{
		EntryId:       uuid.MustParse(expid),
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

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Create(&expense).Error

	// TODO: nicer error handling
	if err != nil {
		common.EmitError(ctx, NewCreateExpenseFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not create Expense: %s", err.Error())))
		return
	}

	resp := casheerapi.CreateExpenseResponse{
		Data: ExpenseToPublic(expense),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
