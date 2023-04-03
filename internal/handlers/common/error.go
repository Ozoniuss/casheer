package common

import (
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
		Error: err,
	})
}
