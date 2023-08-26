package expenses

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) HandleListExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	params, err := common.CtxGetTyped[casheerapi.ListExpenseParams](ctx, "queryparams")
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	var expenses []model.Expense
	db := h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).Where("entry_id = ?", entid).Session(&gorm.Session{})
	db = applyFilters(db, params)

	err = db.Debug().Order("amount desc").Order("created_at desc").Find(&expenses).Error
	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.ListExpenseResponse{
		Data: make([]casheerapi.ExpenseListItemData, 0, len(expenses)),
	}
	for _, e := range expenses {
		resp.Data = append(resp.Data, ExpenseToPublicList(e, h.entriesURL))
	}
	ctx.JSON(http.StatusOK, &resp)
}

func applyFilters(db *gorm.DB, params casheerapi.ListExpenseParams) *gorm.DB {
	if params.Currency != nil {
		db = db.Where("currency = ?", *params.Currency)
	} else {
		db = db.Order("currency")
	}
	if params.Name != nil {
		db = db.Where("name = ?", *params.Name)
	}
	if params.PaymentMethod != nil {
		db = db.Where("payment_method = ?", *params.PaymentMethod)
	}
	if params.AmountGt != nil {
		db = db.Where("amount >= ?", *params.AmountGt)
	}
	if params.AmountLt != nil {
		db = db.Where("amount <= ?", *params.AmountLt)
	}
	return db
}
