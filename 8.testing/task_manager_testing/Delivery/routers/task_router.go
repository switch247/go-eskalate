package router

import (
	"main/Delivery/controllers"
	"main/Repositories"
	"main/UseCases"
	"main/config"

	"main/Infrastructure"

	"github.com/gin-gonic/gin"
)

func TaskRouter(r *gin.Engine) error {
	taskRouter := r.Group("/tasks", Infrastructure.AuthMiddleware())
	{

		client, err := config.GetClient()
		if err != nil {
			return err
		}
		DataBase := client.Database("test")
		collection := DataBase.Collection("tasks")

		taskRepository, err := Repositories.NewTaskRepository(client, DataBase, collection)
		taskUseCase, err := UseCases.NewTaskUseCase(taskRepository)
		taskController, err := controllers.NewTaskController(taskUseCase)
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
