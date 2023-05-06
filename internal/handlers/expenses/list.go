package expenses

import (
	"fmt"
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) HandleListExpense(ctx *gin.Context) {

	// At the moment only listing the expenses for a certain entry; it doesn't
	// really make sense to list all of them.
	entid := ctx.GetString("entid")

	var params casheerapi.ListExpenseParams
	err := ctx.ShouldBindQuery(&params)

	if err != nil {
		common.EmitError(ctx, NewListExpenseFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	// This is needed to query using zero values as well, see
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	var filters = make(map[string]any)

	if params.Value != nil {
		filters["value"] = *params.Value
	}
	if params.Description != nil {
		filters["description"] = *params.Description
	}
	if params.Name != nil {
		filters["name"] = *params.Name
	}
	if params.PaymentMethod != nil {
		filters["payment_method"] = *params.PaymentMethod
	}

	var expenses []model.Expense
	err = h.db.WithContext(ctx).Scopes(model.RequiredEntry(uuid.MustParse(entid))).
		Where(filters).Order("value desc").Order("created_at desc").Find(&expenses).Error

	if err != nil {
		common.EmitError(ctx, NewListExpenseFailedError(
			http.StatusInternalServerError,
			fmt.Sprintf("Could not list entries: %s", err.Error())))
		return
	}

	resp := casheerapi.ListExpenseResponse{
		Data: make([]casheerapi.ExpenseListItemData, 0, len(expenses)),
		Links: casheerapi.ListExpenseLinks{
			Self:    fmt.Sprintf("%s/%s/expenses/", h.apiPaths.Entries, entid),
			Entries: h.apiPaths.Entries + "/",
			Debts:   h.apiPaths.Debts + "/",
		},
	}
	for _, e := range expenses {
		resp.Data = append(resp.Data, ExpenseToPublicList(e, h.apiPaths))
	}
	ctx.JSON(http.StatusOK, &resp)
}
