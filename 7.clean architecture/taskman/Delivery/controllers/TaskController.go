package controllers

import (
	"fmt"

	"main/Repositories"

	"main/Domain"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController interface {
	GetAllTasks(c *gin.Context)
	CreateTasks(c *gin.Context)
	GetTasksById(c *gin.Context)
	UpdateTasksById(c *gin.Context)
	DeleteTasksById(c *gin.Context)
}

type taskController struct {
	TaskRepository Domain.TaskRepository
}

func NewTaskController() (*taskController, error) {
	service_reference, err := Repositories.NewTaskRepository()
	if err != nil {
		return nil, err
	}
	return &taskController{
		TaskRepository: service_reference,
	}, nil
}

func (tc *taskController) GetAllTasks(c *gin.Context) {
	// Get the user ID from the context
	logedUser, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}

	// var newTak Domain.Task
	// v := validator.New()
	// fmt.Print(v, newTak)
	tasks, err, statusCode := tc.TaskRepository.GetTasks(logedUser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)

}

func (tc *taskController) CreateTasks(c *gin.Context) {
	// Retrieve the user information from the context
	// Get the user ID from the context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to retrieve user ID"})
		c.Abort()
		return
	}
	fmt.Printf("User: %v:", userID)
	var newTak Domain.Task
	v := validator.New()
	if err := c.ShouldBindJSON(&newTak); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}
	if err := v.Struct(newTak); err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		return
	}
	err := utils.ValidateStatus(&newTak)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data, allowed options are Pending, In Progress, and Completed", "error": err.Error()})
		return
	}

	newTak.User_ID = userID.(string)

	//
	task, err, statusCode := tc.TaskRepository.CreateTasks(&newTak)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, task)
	}
}

func (tc *taskController) GetTasksById(c *gin.Context) {

	logedUser, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err, statusCode := tc.TaskRepository.GetTasksById(objectID, logedUser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(200, task)
	}
}

func (tc *taskController) UpdateTasksById(c *gin.Context) {
	logedUser, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedTask Domain.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//validate status first
	er := utils.ValidateStatus(&updatedTask)
	if er != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data, allowed options are Pending, In Progress, and Completed", "error": err.Error()})
		return
	}

	data, er, status := tc.TaskRepository.UpdateTasksById(objectID, updatedTask, logedUser)
	if er != nil {
		c.IndentedJSON(status, gin.H{"error": er.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Task updated", "data": data})
	}

}

func (tc *taskController) DeleteTasksById(c *gin.Context) {
	logedUser, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, statusCode := tc.TaskRepository.DeleteTasksById(objectID, logedUser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"success": "Data Deleted"})
	}
}
