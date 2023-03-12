package debts

import (
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandleListDebt(ctx *gin.Context) {

	var params casheerapi.ListDebtParams
	err := ctx.BindQuery(&params)

	if err != nil {
		common.EmitError(ctx, NewListDebtFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	if params.Person != nil {
		filters["person"] = params.Person
	}

	var debts []model.Debt
	err = h.db.WithContext(ctx).Where(filters).Order("person asc").Order("amount desc").Order("uuid asc").Find(&debts).Error

	if err != nil {
		common.EmitError(ctx, NewListDebtFailedError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list debts: %s", err.Error())))
		return
	}

	resp := casheerapi.ListDebtResponse{
		Data: make([]casheerapi.DebtData, 0, len(debts)),
	}
	for _, e := range debts {
		resp.Data = append(resp.Data, DebtToPublic(e))
	}
	ctx.JSON(http.StatusOK, &resp)
}
