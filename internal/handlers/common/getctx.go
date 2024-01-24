package common

import (
	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// CtxGetTypes is very similar to ctx.Get(), except it makes use of generics.
func CtxGetTyped[
	T casheerapi.CreateEntryRequest |
		casheerapi.CreateExpenseRequest |
		casheerapi.CreateDebtRequest |
		casheerapi.UpdateEntryRequest |
		casheerapi.UpdateExpenseRequest |
		casheerapi.UpdateDebtRequest |
		casheerapi.ListExpenseParams |
		casheerapi.ListEntryParams |
		casheerapi.ListDebtParams |
		casheerapi.GetEntryParams](ctx *gin.Context, param string) (T, error) {
	var req T
	reqval, ok := ctx.Get(param)
	if !ok {
		return req, ierrors.NewMissingGinContextParamError(param)
	}

	req, ok = reqval.(T)
	if !ok {
		return req, ierrors.NewInvalidContextParamTypeError(param)
	}

	return req, nil
}
