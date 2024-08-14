package UseCases_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/Domain"
	"main/UseCases"
	"main/mocks"
	"main/utils"
)

type TaskUseCaseTestSuite struct {
	suite.Suite
	mockTaskRepo *mocks.TaskRepository
	taskUseCase  Domain.TaskUseCase
	ctx          *gin.Context
}

func (suite *TaskUseCaseTestSuite) SetupTest() {
	suite.mockTaskRepo = new(mocks.TaskRepository)
	var err error
	suite.ctx = &gin.Context{}
	suite.taskUseCase, err = UseCases.NewTaskUseCase(suite.mockTaskRepo)
	assert.NoError(suite.T(), err)
}

func (suite *TaskUseCaseTestSuite) TestGetAllTasks() {
	// Arrange
	obj_id := primitive.NewObjectID()
	loggedUser := Domain.OmitedUser{
		ID:    obj_id,
		Email: "test@example.com",
	}

	tasks := []*Domain.Task{
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 1",
			Description: "Description 1",
			User_ID:     utils.ObjectIdToString(loggedUser.ID),
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "Task 2",
			Description: "Description 2",
			User_ID:     utils.ObjectIdToString(loggedUser.ID),
		},
	}

	suite.mockTaskRepo.On("GetTasks", mock.Anything, loggedUser).Return(tasks, nil, 200).Once()

	// Act
	result, err, statusCode := suite.taskUseCase.GetAllTasks(suite.ctx, loggedUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), tasks, result)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestCreateTasks() {
	// Arrange
	obj_id := primitive.NewObjectID()
	task := &Domain.Task{
		Title:       "New Task",
		Description: "Description",
		User_ID:     utils.ObjectIdToString(obj_id),
	}

	createdTask := Domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       "New Task",
		Description: "Description",
		User_ID:     utils.ObjectIdToString(obj_id),
	}

	suite.mockTaskRepo.On("CreateTasks", mock.Anything, task).Return(createdTask, nil, 201).Once()

	// Act
	result, err, statusCode := suite.taskUseCase.CreateTasks(suite.ctx, task)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, statusCode)
	assert.Equal(suite.T(), createdTask, result)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestGetTasksById() {
	// Arrange
	obj_id := primitive.NewObjectID()
	loggedUser := Domain.OmitedUser{
		ID:    obj_id,
		Email: "test@example.com",
	}

	taskID := primitive.NewObjectID()
	task := Domain.Task{
		ID:          taskID,
		Title:       "Task 1",
		Description: "Description 1",
		User_ID:     utils.ObjectIdToString(loggedUser.ID),
	}

	suite.mockTaskRepo.On("GetTasksById", mock.Anything, taskID, loggedUser).Return(task, nil, 200).Once()

	// Act
	result, err, statusCode := suite.taskUseCase.GetTasksById(suite.ctx, taskID, loggedUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), task, result)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestUpdateTasksById() {
	// Arrange
	obj_id := primitive.NewObjectID()
	loggedUser := Domain.OmitedUser{
		ID:    obj_id,
		Email: "test@example.com",
	}

	taskID := primitive.NewObjectID()
	updatedTask := Domain.Task{
		Title:       "Updated Task",
		Description: "Updated Description",
		User_ID:     utils.ObjectIdToString(loggedUser.ID),
	}

	updatedTaskResult := Domain.Task{
		ID:          taskID,
		Title:       "Updated Task",
		Description: "Updated Description",
		User_ID:     utils.ObjectIdToString(loggedUser.ID),
	}

	suite.mockTaskRepo.On("UpdateTasksById", mock.Anything, taskID, updatedTask, loggedUser).Return(updatedTaskResult, nil, 200).Once()

	// Act
	result, err, statusCode := suite.taskUseCase.UpdateTasksById(suite.ctx, taskID, updatedTask, loggedUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), updatedTaskResult, result)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestDeleteTasksById() {
	// Arrange
	obj_id := primitive.NewObjectID()
	loggedUser := Domain.OmitedUser{
		ID:    obj_id,
		Email: "test@example.com",
	}

	taskID := primitive.NewObjectID()

	suite.mockTaskRepo.On("DeleteTasksById", mock.Anything, taskID, loggedUser).Return(nil, 200).Once()

	// Act
	err, statusCode := suite.taskUseCase.DeleteTasksById(suite.ctx, taskID, loggedUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

func (suite *TaskUseCaseTestSuite) TestGetAllTasksError() {
	// Arrange
	obj_id := primitive.NewObjectID()
	loggedUser := Domain.OmitedUser{
		ID:    obj_id,
		Email: "test@example.com",
	}

	suite.mockTaskRepo.On("GetTasks", mock.Anything, loggedUser).Return(nil, errors.New("get tasks error"), 500).Once()

	// Act
	result, err, statusCode := suite.taskUseCase.GetAllTasks(suite.ctx, loggedUser)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 500, statusCode)
	assert.Nil(suite.T(), result)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

func TestTaskUseCase(t *testing.T) {

	suite.Run(t, new(TaskUseCaseTestSuite))
}
