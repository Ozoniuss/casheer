package middlewares

import (
	"errors"

	"github.com/Ozoniuss/casheer/internal/apierrors"
	ierrors "github.com/Ozoniuss/casheer/internal/errors"
	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/handlers/debts"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/gin-gonic/gin"
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

		var missingContextParamError ierrors.MissingGinContextParam
		var invalidContextParamTypeError ierrors.InvalidGinContextParamType
		switch {

		case errors.As(err, &missingContextParamError):
			common.EmitError(ctx, apierrors.NewUnknownError("something went wrong while handling the request."))
			return

		case errors.As(err, &invalidContextParamTypeError):
			common.EmitError(ctx, apierrors.NewUnknownError("something went wrong while handling the request."))
			return

		case errors.Is(err, model.InvalidDebtErr{}):
			common.EmitError(ctx,
				debts.NewInvalidDebtError(err.Error()))
			return

		case errors.Is(err, gorm.ErrRecordNotFound):
			common.EmitError(ctx, apierrors.NewNotFoundError())
			return

		default:
			common.EmitError(ctx, apierrors.NewUnknownError(err.Error()))
		}

	}
}
