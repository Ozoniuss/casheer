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

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	var updatedFields = make(map[string]any)

	if req.Person != nil {
		updatedFields["person"] = *req.Person
	}
	if req.Amount != nil {
		updatedFields["amount"] = *req.Amount
	}
	if req.Details != nil {
		updatedFields["details"] = *req.Details
	}

	var Debt model.Debt
	err = h.db.WithContext(ctx).Model(&Debt).Clauses(clause.Returning{}).Where("id = ?", id).Updates(updatedFields).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.UpdateDebtResponse{
		Data: DebtToPublic(Debt, h.debtsURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}
