package router

import (
	"main/Delivery/controllers"
	"main/config"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine) error {
	authRouter := r.Group("/auth")
	{
		client, err := config.GetClient()
		if err != nil {
			return err
		}
		DataBase := client.Database("test")
		_collection := DataBase.Collection("users")
		authController, err := controllers.NewAuthController(client, DataBase, _collection)
		if err != nil {
			return err
		}
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/register", authController.Register)
	}
	return nil
}
