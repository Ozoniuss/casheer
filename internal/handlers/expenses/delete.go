package expenses

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *handler) HandleDeleteExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	id := ctx.GetInt("expid")

	expense := model.Expense{}
	err := h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).Clauses(clause.Returning{}).
		Where("id = ?", id).Delete(&expense).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
		return
	}

	resp := casheerapi.CreateExpenseResponse{
		Data: ExpenseToPublic(expense, h.entriesURL),
	}

	resp.Data.Links.Self = ""
	ctx.JSON(http.StatusOK, &resp)
}
