package expenses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/Ozoniuss/casheer/internal/currency"
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

	value, err := currency.NewValueBasedOnCurrency(req.Amount, req.Currency, req.Exponent)
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	expense := model.Expense{
		EntryId:       entid,
		Value:         value,
		Description:   req.Description,
		Name:          req.Name,
		PaymentMethod: req.PaymentMethod,
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
