package router

import (
	"main/Delivery/controllers"

	"main/Infrastructure"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) error {
	userRouter := r.Group("/users", Infrastructure.AuthMiddleware())
	{
		userController, err := controllers.NewUserController()
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
