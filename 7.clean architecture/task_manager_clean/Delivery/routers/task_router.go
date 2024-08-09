package router

import (
	"main/Delivery/controllers"

	"main/Infrastructure"

	"github.com/gin-gonic/gin"
)

func TaskRouter(r *gin.Engine) error {
	taskRouter := r.Group("/tasks", Infrastructure.AuthMiddleware())
	{
		taskController, err := controllers.NewTaskController()
		if err != nil {
			return err
		}
		taskRouter.GET("", taskController.GetAllTasks)
		taskRouter.POST("", taskController.CreateTasks)
		taskRouter.GET("/:id", taskController.GetTasksById)
		taskRouter.PUT("/:id", taskController.UpdateTasksById)
		taskRouter.DELETE("/:id", taskController.DeleteTasksById)
	}
	return nil
}
