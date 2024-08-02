package controllers

import (
	"fmt"
	"main/data"
	"main/models"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type TaskController interface {
	GetAllTasks(c *gin.Context)
	CreateTasks(c *gin.Context)
	GetTasksById(c *gin.Context)
	UpdateTasksById(c *gin.Context)
	DeleteTasksById(c *gin.Context)
}

type taskController struct {
	taskService data.TaskService
}

func NewTaskController() (*taskController, error) {
	service_reference, err := data.NewTaskService()
	if err != nil {
		return nil, err
	}
	return &taskController{
		taskService: *service_reference,
	}, nil
}

func (tc *taskController) GetAllTasks(c *gin.Context) {
	// var newTak models.Task
	// v := validator.New()
	// fmt.Print(v, newTak)
	tasks, err, statusCode := tc.taskService.GetTasks()
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)

}

func (tc *taskController) CreateTasks(c *gin.Context) {
	var newTak models.Task
	v := validator.New()
	if err := c.ShouldBindJSON(&newTak); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}
	if err := v.Struct(newTak); err != nil {
		fmt.Printf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		return
	}
	err := utils.ValidateStatus(&newTak)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
	}

	//
	task, err, statusCode := tc.taskService.CreateTasks(&newTak)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, task)
	}
}

func (tc *taskController) GetTasksById(c *gin.Context) {
	id := c.Param("id")

	task, err, statusCode := tc.taskService.GetTasksById(id)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(200, task)
	}
}

func (tc *taskController) UpdateTasksById(c *gin.Context) {
	id := c.Param("id")

	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//
	data, er := tc.taskService.UpdateTasksById(id, updatedTask)
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Task updated", "data": data})
	}

}

func (tc *taskController) DeleteTasksById(c *gin.Context) {
	id := c.Param("id")

	err, statusCode := tc.taskService.DeleteTasksById(id)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Data Deleted"})
	}
}
