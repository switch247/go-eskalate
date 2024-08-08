package Domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "Pending"
	StatusCompleted  TaskStatus = "Completed"
	StatusInProgress TaskStatus = "In Progress"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description" validate:"required"`
	DueDate     time.Time          `json:"due_date" validate:"required"`
	Status      TaskStatus         `json:"status" validate:"required"`
	User_ID     string             `json:"user_id"`
}

// oneof='Pending' 'Completed' 'In Progress'

type TaskRepository interface {
	CreateTasks(ctx context.Context, task *Task) (Task, error, int)
	GetTasks(ctx context.Context, user OmitedUser) ([]*Task, error, int)
	GetTasksById(ctx context.Context, id primitive.ObjectID, user OmitedUser) (Task, error, int)
	UpdateTasksById(ctx context.Context, id primitive.ObjectID, task Task, user OmitedUser) (Task, error, int)
	DeleteTasksById(ctx context.Context, id primitive.ObjectID, user OmitedUser) (error, int)
}

type TaskUseCase interface {
	GetAllTasks(c *gin.Context, loggedUser OmitedUser) ([]*Task, error, int)
	CreateTasks(c *gin.Context, task *Task) (Task, error, int)
	GetTasksById(c *gin.Context, id primitive.ObjectID, loggedUser OmitedUser) (Task, error, int)
	UpdateTasksById(c *gin.Context, id primitive.ObjectID, task Task, loggedUser OmitedUser) (Task, error, int)
	DeleteTasksById(c *gin.Context, id primitive.ObjectID, user OmitedUser) (error, int)
}
