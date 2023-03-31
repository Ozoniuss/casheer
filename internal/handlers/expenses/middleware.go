package expenses

import (
	"fmt"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequiredRxpenseUUID is a middleware that ensures the entry uuid is provided
// in the URI.
func RequiredEntryUUID() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		expid := ctx.Param("entid")
		if expid == "" {
			common.EmitError(ctx, NewInvalidEntryError(
				fmt.Sprintf("No uuid provided for entry."),
			))
			ctx.Abort()
			return
		}

		uuid, err := uuid.Parse(expid)
		if err != nil {
			common.EmitError(ctx, NewInvalidEntryError(
				fmt.Sprintf("Invalid uuid format for entry: %s", expid),
			))
			ctx.Abort()
			return
		}

		ctx.Set("entid", uuid.String())
		ctx.Next()
	}

}
