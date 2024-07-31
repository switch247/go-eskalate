package main

import (
	"fmt"
	"net/http"
	"time"

	"main/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Task struct {
	ID          string    `json:"id"  validate:"required"`
	Title       string    `json:"title"  validate:"required""`
	Description string    `json:"description"  validate:"required"`
	DueDate     time.Time `json:"due_date"  validate:"required"`
	Status      string    `json:"status"  validate:"required"`
}

// Mock data for tasks
var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func main() {

	fmt.Println("Task Manager API")
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/test", testAnyjson)

	router.POST("/tasks", createTasks)
	router.GET("/tasks", getAllTasks)
	router.GET("/tasks/:id", getTasksById)
	router.PUT("/tasks/:id", updateTasksById)
	router.DELETE("/tasks/:id", deleteTasksById)

	port := 8080
	test := "localhost:" + string(port)
	fmt.Println(test)
	router.Run()
	// router.Run()localhost:8080")
}

func createTasks(c *gin.Context) {
	var newTak Task
	v := validator.New()

	if err := c.ShouldBindJSON(&newTak); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing data"})
		return
	}
	if err := v.Struct(newTak); err != nil {
		fmt.Printf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	tasks = append(tasks, newTak)
	c.IndentedJSON(http.StatusCreated, tasks)
}
func getAllTasks(c *gin.Context) {
	c.IndentedJSON(200, tasks)
}

func getTasksById(c *gin.Context) {
	id := c.Param("id")
	for _, val := range tasks {
		if val.ID == id {

			c.IndentedJSON(200, val)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})

}

func updateTasksById(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			// Update only the specified fields
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			if !updatedTask.DueDate.IsZero() {
				tasks[i].DueDate = updatedTask.DueDate
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})

}

func deleteTasksById(c *gin.Context) {
	id := c.Param("id")
	v := validator.New()
	var oldTaks Task
	var oldIndex int
	for idx, val := range tasks {
		if val.ID == id {
			oldTaks = val
			oldIndex = idx
		}
	}

	if err := v.Struct(oldTaks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	tasks = append(tasks[:oldIndex], tasks[oldIndex+1:]...)
	c.JSON(http.StatusOK, gin.H{"success": "Data Deleted"})
}

func testAnyjson(ctx *gin.Context) {
	val, err := utils.ReadJSON(ctx)
	if err != nil {
		fmt.Println("failed")
		fmt.Println(err.Error())
		ctx.Status(http.StatusBadRequest)

	} else {
		fmt.Println(val["id"])
		if val["xid"] != nil {
			ctx.JSON(200, val["xid"])
		}
		return
	}
}
