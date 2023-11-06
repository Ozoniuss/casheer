package expenses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/Ozoniuss/casheer/currency"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func (h *handler) HandleCreateExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	req, err := common.CtxGetTyped[casheerapi.CreateExpenseRequest](ctx, "req")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	value, err := currency.NewValueBasedOnCurrency(req.Data.Attributes.Value.Amount, req.Data.Attributes.Value.Currency, req.Data.Attributes.Value.Exponent)
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	expense := model.Expense{
		EntryId: entid,
		Value:   model.FromCurrencyValue(value),
		Name:    req.Data.Attributes.Name,
	}

	if req.Data.Attributes.Description != nil {
		expense.Description = *req.Data.Attributes.Description
	}
	if req.Data.Attributes.PaymentMethod != nil {
		expense.PaymentMethod = *req.Data.Attributes.PaymentMethod
	}

	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(expense.EntryId)).Clauses(clause.Returning{}).Create(&expense).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.CreateExpenseResponse{
		Data: ExpenseToPublic(expense, h.entriesURL),
	}

	ctx.JSON(http.StatusCreated, &resp)
}
