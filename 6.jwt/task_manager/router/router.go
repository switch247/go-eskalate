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
		taskRouter.GET("", taskController.GetAllTasks)
		taskRouter.POST("", taskController.CreateTasks)
		taskRouter.GET("/:id", taskController.GetTasksById)
		taskRouter.PUT("/:id", taskController.UpdateTasksById)
		taskRouter.DELETE("/:id", taskController.DeleteTasksById)
	}
	userRouter := r.Group("/users", middleware.AuthMiddleware())
	{
		userController, err := controllers.NewUserController()
		if err != nil {
			return nil
		}
		userRouter.POST("", middleware.IsAdminMiddleware(), userController.CreateUser)
		userRouter.GET("", middleware.IsAdminMiddleware(), userController.GetUsers)
		userRouter.GET("/:id", userController.GetUser)
		userRouter.PUT("/:id", userController.UpdateUser)
		userRouter.DELETE("/:id", userController.DeleteUser)
	}

	return r
}
