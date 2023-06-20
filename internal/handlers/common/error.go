package common

import (
	"fmt"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// EmitError sends an error response back to the client.
func EmitError(ctx *gin.Context, err public.Error) {
	// ctx.JSON(err.Status, gin.H{
	// 	"error": err,
	// })
	ctx.JSON(err.Status, casheerapi.ErrorResponse{
		Err: err,
	})
}

// CtxGetTypes is very similar to ctx.Get(), except it makes use of generics.
func CtxGetTyped[T casheerapi.CreateEntryRequest](ctx *gin.Context, param string) (T, bool) {
	reqval, ok := ctx.Get(param)
	if !ok {
		EmitError(ctx, ierrors.NewMissingContextParamError(
			fmt.Sprintf("Context parameter %s not found.", param),
		))
		return T{}, false
	}

	req, ok := reqval.(T)
	if !ok {
		EmitError(ctx, ierrors.NewInvalidContextParamError(
			fmt.Sprintf("Context parameter %s has invalid type.", param),
		))
		return T{}, false
	}

	return req, true
}
