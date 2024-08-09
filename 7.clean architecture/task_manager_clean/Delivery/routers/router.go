package router

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}
func NewRouter() *gin.Engine {

	// Create a sub-router group for authentication routes
	AuthRouter(router)
	// Create a sub-router group for task routes
	TaskRouter(router)
	// Create a sub-router group for user routes
	UserRouter(router)

	return router
}
