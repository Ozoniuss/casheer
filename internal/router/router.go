package router

import (
	"github.com/gin-gonic/gin"
)

// NewRouter initializes the gin router with the existing handlers and options.
func NewRouter() (*gin.Engine, error) {
	r := gin.Default()
	return r, nil
}
