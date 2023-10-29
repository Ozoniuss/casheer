package expenses

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleUpdateExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	id := ctx.GetInt("expid")

	req, err := common.CtxGetTyped[casheerapi.UpdateExpenseRequest](ctx, "req")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	var oldExpense model.Expense
	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).Where("id = ?", id).Take(&oldExpense).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}
	updateExpenseFields(req, &oldExpense)

	err = h.db.WithContext(ctx).Clauses(clause.Returning{}).Scopes(model.ValidateModelScope[model.Expense](oldExpense)).
		Save(&oldExpense).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
	}
	resp := casheerapi.UpdateExpenseResponse{
		Data: ExpenseToPublic(oldExpense, h.entriesURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}

func updateExpenseFields(req casheerapi.UpdateExpenseRequest, expense *model.Expense) {
	if req.Data.Attributes.Amount != nil {
		expense.Amount = *req.Data.Attributes.Amount
	}
	if req.Data.Attributes.Currency != nil {
		expense.Currency = *req.Data.Attributes.Currency
	}
	if req.Data.Attributes.Description != nil {
		expense.Description = *req.Data.Attributes.Description
	}
	if req.Data.Attributes.Exponent != nil {
		expense.Exponent = *req.Data.Attributes.Exponent
	}
	if req.Data.Attributes.Name != nil {
		expense.Name = *req.Data.Attributes.Name
	}
	if req.Data.Attributes.PaymentMethod != nil {
		expense.PaymentMethod = *req.Data.Attributes.PaymentMethod
	}
}

func getUpdatedFields(req casheerapi.UpdateExpenseRequest) (model.Expense, []string) {

	// See https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make([]string, 0, 6)
	expense := model.Expense{}

	// TODO: proper validation here.
	if req.Data.Attributes.Amount != nil {
		updatedFields = append(updatedFields, "amount")
		expense.Amount = *req.Data.Attributes.Amount
	}
	if req.Data.Attributes.Currency != nil {
		updatedFields = append(updatedFields, "currency")
		expense.Currency = *req.Data.Attributes.Currency
	}
	if req.Data.Attributes.Description != nil {
		updatedFields = append(updatedFields, "description")
		expense.Description = *req.Data.Attributes.Description
	}
	if req.Data.Attributes.Exponent != nil {
		updatedFields = append(updatedFields, "exponent")
		expense.Exponent = *req.Data.Attributes.Exponent
	}
	if req.Data.Attributes.Name != nil {
		updatedFields = append(updatedFields, "name")
		expense.Name = *req.Data.Attributes.Name
	}
	if req.Data.Attributes.PaymentMethod != nil {
		updatedFields = append(updatedFields, "payment_method")
		expense.PaymentMethod = *req.Data.Attributes.PaymentMethod
	}

	return expense, updatedFields
}
