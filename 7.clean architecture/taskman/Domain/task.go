package Domain

import (
	"time"

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
	GetTasks(user OmitedUser) ([]*Task, error, int)
	CreateTasks(task *Task) (Task, error, int)
	GetTasksById(id primitive.ObjectID, user OmitedUser) (Task, error, int)
	UpdateTasksById(id primitive.ObjectID, task Task, user OmitedUser) (Task, error, int)
	DeleteTasksById(id primitive.ObjectID, user OmitedUser) (error, int)
}

type TaskUsecase interface {
	GetTasks(user OmitedUser) ([]*Task, error, int)
	CreateTasks(task *Task) (Task, error, int)
	GetTasksById(id primitive.ObjectID, user OmitedUser) (Task, error, int)
	UpdateTasksById(id primitive.ObjectID, task Task, user OmitedUser) (Task, error, int)
	DeleteTasksById(id primitive.ObjectID, user OmitedUser) (error, int)
}
