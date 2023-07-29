package common

import (
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
)

// EmitError sends an error response back to the client.
func EmitError(ctx *gin.Context, err public.Error) {
	ctx.JSON(err.Status, casheerapi.ErrorResponse{
		Err: err,
	})
	ctx.Abort()
}

// ErrorAndAbort adds an error to the context and aborts the pending handlers.
// This method is surprisingly not part of gin.
func ErrorAndAbort(ctx *gin.Context, err error) {
	ctx.Error(err)
	ctx.Abort()
}
