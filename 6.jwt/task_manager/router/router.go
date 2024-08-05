package router

import (
	"main/controllers"

	"main/middleware"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	r = gin.Default()
}
func NewRouter() *gin.Engine {

	// Create a sub-router group for authentication routes
	authRouter := r.Group("/auth")
	{
		authController, err := controllers.NewAuthController()
		if err != nil {
			return nil
		}
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/register", authController.Register)
	}

	// Create a sub-router group for task routes
	taskRouter := r.Group("/tasks", middleware.AuthMiddleware())
	{
		taskController, err := controllers.NewTaskController()
		if err != nil {
			return nil
		}
		taskRouter.GET("", middleware.UserMiddleware(), taskController.GetAllTasks)
		taskRouter.POST("", middleware.UserMiddleware(), taskController.CreateTasks)
		taskRouter.GET("/:id", taskController.GetTasksById)
		taskRouter.PUT("/:id", middleware.UserMiddleware(), taskController.UpdateTasksById)
		taskRouter.DELETE("/:id", middleware.UserMiddleware(), taskController.DeleteTasksById)
	}
	userRouter := r.Group("/users", middleware.AuthMiddleware())
	{
		userController, err := controllers.NewUserController()
		if err != nil {
			return nil
		}
		userRouter.GET("", userController.GetUsers)
		userRouter.GET("/:id", userController.GetUser)
		userRouter.DELETE("/:id", userController.DeleteUser)
	}

	return r
}