package models

import "time"

type Task struct {
	ID          string    `json:"id"  validate:"required"`
	Title       string    `json:"title"  validate:"required""`
	Description string    `json:"description"  validate:"required"`
	DueDate     time.Time `json:"due_date"  validate:"required"`
	Status      string    `json:"status"  validate:"required"`
}
