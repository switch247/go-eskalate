package router

import (
	"main/Delivery/controllers"
	"main/Repositories"
	"main/UseCases"
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
		authRepository, err := Repositories.NewAuthRepository(client, DataBase, _collection)
		authUseCase, err := UseCases.NewAuthUseCase(authRepository)
		authController, err := controllers.NewAuthController(authUseCase)
		if err != nil {
			return err
		}
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/register", authController.Register)
	}
	return nil
}
