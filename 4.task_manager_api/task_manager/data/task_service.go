package data

import (
	// "errors"
	// "strconv"

	"time"

	"main/models"

	"errors"

	"github.com/go-playground/validator"
)

type taskService interface {
	GetTasks() ([]models.Task, error)
	CreateTasks(task *models.Task) (models.Task, error)
	GetTasksById(id string) (models.Task, error)
	UpdateTasksById(id string, task models.Task) (models.Task, error)
	DeleteTasksById(id string) error
}

type TaskService struct {
	validator *validator.Validate
	tasks     map[string]*models.Task
	// mu        sync.RWMutex //i will add this back once i understand routines properly
}

func NewTaskService() *TaskService {

	ts := &TaskService{
		validator: validator.New(),
		tasks:     make(map[string]*models.Task),
	}
	ts.tasks["1"] = &models.Task{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"}
	ts.tasks["2"] = &models.Task{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"}
	ts.tasks["3"] = &models.Task{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"}
	return ts
}

func (ts *TaskService) GetTasks() ([]models.Task, error) {
	// ts.mu.RLock()
	// defer ts.mu.RUnlock()

	var tasks []models.Task
	for _, task := range ts.tasks {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (ts *TaskService) CreateTasks(task *models.Task) (models.Task, error) {
	_, err := ts.tasks[task.ID]
	if err == true {
		return models.Task{}, errors.New("Task already exists")
	} else {
		ts.tasks[task.ID] = task
	}
	return *task, nil
}

func (ts *TaskService) GetTasksById(id string) (models.Task, error) {
	ans, ok := ts.tasks[id]
	if ok == true {
		return *ans, nil
	}
	return models.Task{}, errors.New("Task not found")
}

func (ts *TaskService) UpdateTasksById(id string, task models.Task) (models.Task, error) {
	updatedTask, ok := ts.tasks[id]
	NewTask := models.Task{ID: id}
	if ok == false {
		return models.Task{}, errors.New("Task does not exists")
	} else {
		// Update only the specified fields
		if updatedTask.Title != "" {
			NewTask.Title = updatedTask.Title
		}
		if updatedTask.Description != "" {
			NewTask.Description = updatedTask.Description
		}
		if updatedTask.Status != "" {
			NewTask.Status = updatedTask.Status
		}
		if !updatedTask.DueDate.IsZero() {
			NewTask.DueDate = updatedTask.DueDate
		}

		ts.tasks[id] = &NewTask

		return NewTask, nil
	}
}

func (ts *TaskService) DeleteTasksById(id string) error {
	_, ok := ts.tasks[id]
	if ok == false {
		return errors.New("Task does not exist")
	} else {
		delete(ts.tasks, id)
		return nil
	}
}
