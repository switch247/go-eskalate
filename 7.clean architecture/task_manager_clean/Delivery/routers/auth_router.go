package router

import (
	"main/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine) error {
	authRouter := r.Group("/auth")
	{
		authController, err := controllers.NewAuthController()
		if err != nil {
			return err
		}
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/register", authController.Register)
	}
	return nil
}
