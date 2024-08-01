package models

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "Pending"
	StatusCompleted TaskStatus = "Completed"
	StatusCancelled TaskStatus = "In Progress"
)

type Task struct {
	ID          string     `json:"id"  validate:"required"`
	Title       string     `json:"title"  validate:"required"`
	Description string     `json:"description"  validate:"required"`
	DueDate     time.Time  `json:"due_date"  validate:"required"`
	Status      TaskStatus `json:"status"  validate:"required"`
}
