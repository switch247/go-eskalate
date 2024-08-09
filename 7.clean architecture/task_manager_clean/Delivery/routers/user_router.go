package router

import (
	"main/Delivery/controllers"
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
		userController, err := controllers.NewUserController(client, DataBase, _collection)
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
