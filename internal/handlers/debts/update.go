package debts

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleUpdateDebt(ctx *gin.Context) {

	id := ctx.GetInt("dbtid")
	req, err := common.CtxGetTyped[casheerapi.UpdateDebtRequest](ctx, "req")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	updatedDebt, fields := getUpdatedFields(req)

	err = h.db.WithContext(ctx).Select(fields).Clauses(clause.Returning{}).Scopes(model.ValidateModelScope[model.Debt](updatedDebt)).Where("id = ?", id).Updates(&updatedDebt).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.UpdateDebtResponse{
		Data: DebtToPublic(updatedDebt, h.debtsURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}

func getUpdatedFields(req casheerapi.UpdateDebtRequest) (model.Debt, []string) {

	debt := model.Debt{}
	var updatedFields = make([]string, 0, 3)

	if req.Person != nil {
		updatedFields = append(updatedFields, "person")
		debt.Person = *req.Person
	}
	if req.Amount != nil {
		updatedFields = append(updatedFields, "amount")
		debt.Amount = *req.Amount
	}
	if req.Details != nil {
		updatedFields = append(updatedFields, "details")
		debt.Details = *req.Details
	}

	return debt, updatedFields
}
