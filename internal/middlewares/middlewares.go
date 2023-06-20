package middlewares

import (
	"fmt"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// BindJSONRequest does the same job as ctx.ShouldBindJSON, while also writing
// a custom error to the response if the binding is not succesful.
func BindJSONRequest[T casheerapi.CreateEntryRequest]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			common.EmitError(ctx, ierrors.NewRequestBindingError(
				fmt.Sprintf("Could not bind request body: %s", err.Error()),
			))
			ctx.Abort()
			return
		}
		ctx.Set("req", req)
		ctx.Next()
	}
}

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
