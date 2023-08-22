package expenses

import (
	"net/http"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

func (h *handler) HandleGetExpense(ctx *gin.Context) {

	entid := ctx.GetInt("entid")
	id := ctx.GetInt("expid")

	var expense model.Expense
	err := h.db.WithContext(ctx).Scopes(model.RequiredEntry(entid)).
		Where("id = ?", id).Take(&expense).Error

	if err != nil {
		common.ErrorAndAbort(ctx, err)
	}

	resp := casheerapi.GetExpenseResponse{
		Data: ExpenseToPublic(expense, h.entriesURL),
	}

	ctx.JSON(http.StatusOK, &resp)
}
