package UseCases

import (
	"context"
	"time"

	"main/Domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskUseCase struct {
	TaskRepository Domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUseCase(service_reference Domain.TaskRepository) (*taskUseCase, error) {
	return &taskUseCase{
		TaskRepository: service_reference,
		contextTimeout: time.Second * 10,
	}, nil
}

func (tu *taskUseCase) GetAllTasks(c *gin.Context, loggedUser Domain.OmitedUser) ([]*Domain.Task, error, int) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.TaskRepository.GetTasks(ctx, loggedUser)

}

func (tu *taskUseCase) CreateTasks(c *gin.Context, task *Domain.Task) (Domain.Task, error, int) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.TaskRepository.CreateTasks(ctx, task)

}

func (tu *taskUseCase) GetTasksById(c *gin.Context, id primitive.ObjectID, loggedUser Domain.OmitedUser) (Domain.Task, error, int) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.TaskRepository.GetTasksById(ctx, id, loggedUser)

}

func (tu *taskUseCase) UpdateTasksById(c *gin.Context, id primitive.ObjectID, task Domain.Task, loggedUser Domain.OmitedUser) (Domain.Task, error, int) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.TaskRepository.UpdateTasksById(ctx, id, task, loggedUser)

}

func (tu *taskUseCase) DeleteTasksById(c *gin.Context, id primitive.ObjectID, user Domain.OmitedUser) (error, int) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	return tu.TaskRepository.DeleteTasksById(ctx, id, user)

}
