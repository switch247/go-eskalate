package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "Pending"
	StatusCompleted TaskStatus = "Completed"
	StatusCancelled TaskStatus = "In Progress"
)

type Task struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description" validate:"required"`
	DueDate     time.Time          `json:"due_date" validate:"required"`
	Status      TaskStatus         `json:"status" validate:"required"`
}
