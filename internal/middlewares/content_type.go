package middlewares

import "github.com/gin-gonic/gin"

func JSONApiContentType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/vnd.api+json")
	}
}
