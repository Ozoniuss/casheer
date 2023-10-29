package debts

import (
	"fmt"
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

	var oldDebt model.Debt
	err = h.db.WithContext(ctx).Where("id = ?", id).Take(&oldDebt).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}
	updateDebtFields(req, &oldDebt)

	err = h.db.WithContext(ctx).Scopes(model.ValidateModelScope[model.Debt](oldDebt)).Clauses(clause.Returning{}).Save(&oldDebt).Error
	if err != nil {
		fmt.Printf("%v %T", err, err)
		common.ErrorAndAbort(ctx, err)
		return
	}
	fmt.Println("jhere")

	resp := casheerapi.UpdateDebtResponse{
		Data: DebtToPublic(oldDebt, h.debtsURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}

func updateDebtFields(req casheerapi.UpdateDebtRequest, debt *model.Debt) {
	if req.Data.Attributes.Person != nil {
		debt.Person = *req.Data.Attributes.Person
	}
	if req.Data.Attributes.Amount != nil {
		debt.Amount = *req.Data.Attributes.Amount
	}
	if req.Data.Attributes.Currency != nil {
		debt.Currency = *req.Data.Attributes.Currency
	}
	if req.Data.Attributes.Exponent != nil {
		debt.Exponent = *req.Data.Attributes.Exponent
	}
	if req.Data.Attributes.Details != nil {
		debt.Details = *req.Data.Attributes.Details
	}
}
