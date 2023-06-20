package common

import (
	"fmt"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// CtxGetTypes is very similar to ctx.Get(), except it makes use of generics.
func CtxGetTyped[T casheerapi.CreateEntryRequest | casheerapi.CreateExpenseRequest | casheerapi.CreateDebtRequest | casheerapi.UpdateEntryRequest](ctx *gin.Context, param string) (T, bool) {
	var req T
	reqval, ok := ctx.Get(param)
	if !ok {
		EmitError(ctx, ierrors.NewMissingContextParamError(
			fmt.Sprintf("Context parameter %s not found.", param),
		))
		return req, false
	}

	req, ok = reqval.(T)
	if !ok {
		EmitError(ctx, ierrors.NewInvalidContextParamError(
			fmt.Sprintf("Context parameter %s has invalid type.", param),
		))
		return req, false
	}

	return req, true
}
