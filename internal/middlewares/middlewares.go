package middlewares

import (
	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/gin-gonic/gin"
)

func RequiredParam[T int](paramName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param, ok := ctx.Get(paramName)
		if !ok {
			common.EmitError(ctx, ierrors.NewMissingParamError(paramName))
			ctx.Abort()
			return
		}
		paramVal, ok := param.(T)
		if !ok {
			common.EmitError(ctx, ierrors.NewInvalidParamType(paramName))
			ctx.Abort()
			return
		}
		ctx.Set(paramName, paramVal)
		ctx.Next()
	}
}
