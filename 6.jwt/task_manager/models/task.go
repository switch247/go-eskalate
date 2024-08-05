package models

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
	UserID      primitive.ObjectID `json:"user_id,omitempty"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description" validate:"required"`
	DueDate     time.Time          `json:"due_date" validate:"required"`
	Status      TaskStatus         `json:"status" validate:"required"`
}

// oneof='Pending' 'Completed' 'In Progress'
