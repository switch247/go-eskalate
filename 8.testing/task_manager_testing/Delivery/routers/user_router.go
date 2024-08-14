package router

import (
	"main/Delivery/controllers"
	"main/Repositories"
	"main/UseCases"
	"main/config"

	"main/Infrastructure"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) error {
	userRouter := r.Group("/users", Infrastructure.AuthMiddleware())
	{
		client, err := config.GetClient()
		if err != nil {
			return err
		}
		DataBase := client.Database("test")
		_collection := DataBase.Collection("users")

		userRepository, err := Repositories.NewUserRepository(client, DataBase, _collection)
		if err != nil {
			return err
		}
		userUseCase, err := UseCases.NewUserUseCase(userRepository)
		if err != nil {
			return err
		}
		userController, err := controllers.NewUserController(userUseCase)
		if err != nil {
			return err
		}
		userRouter.POST("", Infrastructure.AuthMiddleware(), userController.CreateUser)
		userRouter.GET("", Infrastructure.IsAdminMiddleware(), userController.GetUsers)
		userRouter.GET("/:id", userController.GetUser)
		userRouter.PUT("/:id", userController.UpdateUser)
		userRouter.DELETE("/:id", userController.DeleteUser)
	}
	return nil
}
