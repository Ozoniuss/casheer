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

	expense, updatedFields := getUpdatedFields(req)
	expense.EntryId = entid
	err = h.db.WithContext(ctx).Select(updatedFields).Scopes(model.RequiredEntry(entid)).Clauses(clause.Returning{}).Scopes(model.ValidateModelScope[model.Expense](expense)).Where("id = ?", id).Updates(&expense).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
	}
	resp := casheerapi.UpdateExpenseResponse{
		Data: ExpenseToPublic(expense, h.entriesURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}

func getUpdatedFields(req casheerapi.UpdateExpenseRequest) (model.Expense, []string) {

	// See https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make([]string, 0, 6)
	expense := model.Expense{}

	// TODO: proper validation here.
	if req.Amount != nil {
		updatedFields = append(updatedFields, "amount")
		expense.Amount = *req.Amount
	}
	if req.Currency != nil {
		updatedFields = append(updatedFields, "currency")
		expense.Currency = *req.Currency
	}
	if req.Description != nil {
		updatedFields = append(updatedFields, "description")
		expense.Description = *req.Description
	}
	if req.Exponent != nil {
		updatedFields = append(updatedFields, "exponent")
		expense.Exponent = *req.Exponent
	}
	if req.Name != nil {
		updatedFields = append(updatedFields, "name")
		expense.Name = *req.Name
	}
	if req.PaymentMethod != nil {
		updatedFields = append(updatedFields, "payment_method")
		expense.PaymentMethod = *req.PaymentMethod
	}

	return expense, updatedFields
}
