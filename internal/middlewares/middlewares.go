package middlewares

import (
	"fmt"
	"strconv"

	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// BindJSONRequest does the same job as ctx.ShouldBindJSON, while also writing
// a custom error to the response if the binding is not succesful.
func BindJSONRequest[
	T casheerapi.CreateEntryRequest |
		casheerapi.CreateDebtRequest |
		casheerapi.CreateExpenseRequest |
		casheerapi.UpdateExpenseRequest |
		casheerapi.UpdateDebtRequest |
		casheerapi.UpdateEntryRequest]() gin.HandlerFunc {
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

func BindQueryParams[
	T casheerapi.ListDebtParams |
		casheerapi.ListEntryParams |
		casheerapi.ListExpenseParams]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params T
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			common.EmitError(ctx, ierrors.NewQueryParamsBindingError(
				fmt.Sprintf("Could not bind query params: %s", err.Error()),
			))
			ctx.Abort()
			return
		}
		ctx.Set("queryparams", params)
		ctx.Next()
	}
}

// GetURLParam does the same job as ctx.Param, while also writing a custom
// error message if the param is not found or is not a valid integer.
//
// For now, all URL params represent identifiers, which are integers.
func GetURLParam(paramName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ok := retrieveAndSetParam(ctx, paramName)
		if !ok {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// GetMultipleURLParam is the same as GetURLParam, except it retrieves values
// for multiple URL parameters.
func GetMultipleURLParam(paramNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, paramName := range paramNames {
			ok := retrieveAndSetParam(ctx, paramName)
			if !ok {
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

func retrieveAndSetParam(ctx *gin.Context, paramName string) bool {
	param := ctx.Param(paramName)
	if param == "" {
		common.EmitError(ctx, ierrors.NewMissingParamError(
			paramName,
		))
		return false
	}
	paramVal, err := strconv.Atoi(param)
	if err != nil {
		common.EmitError(ctx, ierrors.NewInvalidParamType(paramName))
		return false
	}
	ctx.Set(paramName, paramVal)
	return true
}
