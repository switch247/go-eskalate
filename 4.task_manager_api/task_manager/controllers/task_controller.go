package controllers

import (
	"fmt"
	"main/data"
	"main/models"
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

func NewTaskController() *taskController {
	return &taskController{
		taskService: *data.NewTaskService(),
	}
}

func (tc *taskController) GetAllTasks(c *gin.Context) {
	// var newTak models.Task
	// v := validator.New()
	// fmt.Print(v, newTak)
	tasks, err := tc.taskService.GetTasks()
	if err != nil {
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
	//
	task, err := tc.taskService.CreateTasks(&newTak)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, task)
	}
}

func (tc *taskController) GetTasksById(c *gin.Context) {
	id := c.Param("id")

	task, err := tc.taskService.GetTasksById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err := tc.taskService.DeleteTasksById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Data Deleted"})
}
