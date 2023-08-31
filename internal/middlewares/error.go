package middlewares

import (
	"errors"
	"fmt"

	"github.com/Ozoniuss/casheer/internal/apierrors"
	"github.com/Ozoniuss/casheer/internal/currency"
	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

// ErrorHandler is a middleware which gets called when errors occur inside the
// handlers.
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Execute the other handlers in the chain beforehand.
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors[0].Err
		fmt.Printf("An error occured: %s\n", err.Error())

		var invalidQueryParamsError ierrors.InvalidQueryParams
		var invalidJsonBodyError ierrors.InvalidJsonBody
		var missingContextParamError ierrors.MissingGinContextParam
		var invalidContextParamTypeError ierrors.InvalidGinContextParamType
		var invalidResourceError ierrors.InvalidModel
		var sqlite3err sqlite3.Error

		var invalidEntryKey model.ErrExpenseInvalidEntryKey

		var invalidCurrencyError currency.ErrInvalidCurrency

		switch {

		case errors.As(err, &invalidQueryParamsError):
			common.EmitError(ctx, apierrors.NewQueryParamsBindingError(fmt.Sprintf("Invalid query parameters: %s", invalidQueryParamsError)))

		case errors.As(err, &invalidJsonBodyError):
			common.EmitError(ctx, apierrors.NewRequestBindingError(fmt.Sprintf("Invalid request body: %s", invalidJsonBodyError)))

		case errors.As(err, &missingContextParamError):
			common.EmitError(ctx, apierrors.NewUnknownError("something went wrong while handling the request."))
			return

		case errors.As(err, &invalidContextParamTypeError):
			common.EmitError(ctx, apierrors.NewUnknownError("something went wrong while handling the request."))
			return

		case errors.As(err, &invalidResourceError):
			common.EmitError(ctx, apierrors.NewInvalidResourceError(invalidResourceError.Error()))
			return

		// Should catch at least unique constraint error.
		case errors.As(err, &sqlite3err):
			if sqlite3err.Code == sqlite3.ErrConstraint && sqlite3err.ExtendedCode == sqlite3.ErrConstraintUnique {
				common.EmitError(ctx, apierrors.NewAlreadyExistsError(sqlite3err.Error()))
			} else {
				fmt.Println("UNHANDLED")
			}
			return

		case errors.As(err, &invalidCurrencyError):
			common.EmitError(ctx, apierrors.NewInvalidResourceError(invalidCurrencyError.Error()))
			return

		case errors.As(err, &invalidEntryKey):
			common.EmitError(ctx, apierrors.NewInvalidResourceError(invalidEntryKey.Error()))
			return

		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, apierrors.NewNotFoundError())
			return

		default:
			common.EmitError(ctx, apierrors.NewUnknownError(err.Error()))
		}

	}
}
