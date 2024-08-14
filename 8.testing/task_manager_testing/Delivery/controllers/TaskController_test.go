package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"main/Delivery/controllers"
	"main/Domain"
	"main/mocks"
	"main/utils"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock function to extract user from context
func mockExtractUser(c *gin.Context) (Domain.OmitedUser, error) {
	obj_id := primitive.NewObjectID()
	return Domain.OmitedUser{
		ID:       obj_id,
		Is_Admin: true,
	}, nil
}

// TaskControllerTestSuite holds the shared setup for all tests
type TaskControllerTestSuite struct {
	suite.Suite
	controller      controllers.TaskController
	mockTaskUseCase *mocks.TaskUseCase
	mockTaskID      primitive.ObjectID
	userID          string
	is_admin        bool
	patch           *monkey.PatchGuard
}

// Setup initializes the necessary elements for testing
func (suite *TaskControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.controller, _ = controllers.NewTaskController(suite.mockTaskUseCase)
	suite.mockTaskUseCase = new(mocks.TaskUseCase)
	suite.mockTaskID = primitive.NewObjectID()
	suite.userID = primitive.NewObjectID().Hex()
	suite.is_admin = true

	// Patch the ExtractUser function
	suite.patch = monkey.Patch(utils.ExtractUser, mockExtractUser)
}

// TearDown unpatches the monkey patches
func (suite *TaskControllerTestSuite) TearDownTest() {
	suite.patch.Unpatch()
}

func (suite *TaskControllerTestSuite) TestCreateTasks() {
	// Arrange
	mockTask := Domain.Task{
		ID:          suite.mockTaskID,
		Title:       "Test Task",
		Description: "Task description",
		DueDate:     time.Now(),
		Status:      "Pending",
		User_ID:     suite.userID,
	}

	suite.mockTaskUseCase.On("CreateTasks", mock.Anything, mock.Anything).Return(mockTask, nil, http.StatusCreated)

	// Act
	reqBody, _ := json.Marshal(mockTask)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	ctx.Set("user_id", suite.userID)
	ctx.Set("is_admin", suite.is_admin)

	suite.controller.CreateTasks(ctx)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetAllTasks() {

	mockTasks := []*Domain.Task{
		{
			Title:       "Test Task 1",
			Description: "Task description 1",
			DueDate:     time.Now(),
			Status:      "Pending",
			User_ID:     "user123",
		},
		{
			Title:       "Test Task 2",
			Description: "Task description 2",
			DueDate:     time.Now(),
			Status:      "In Progress",
			User_ID:     "user123",
		},
	}

	suite.mockTaskUseCase.On("GetAllTasks", mock.Anything, mock.Anything).Return(mockTasks, nil, http.StatusOK)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	ctx.Set("user_id", suite.userID)
	ctx.Set("is_admin", suite.is_admin)

	suite.controller.GetAllTasks(ctx)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestGetTasksById() {

	mockTask := Domain.Task{
		Title:       "Test Task",
		Description: "Task description",
		DueDate:     time.Now(),
		Status:      "Completed",
		User_ID:     "user123",
	}

	suite.mockTaskUseCase.On("GetTasksById", mock.Anything, suite.mockTaskID, mock.Anything).Return(mockTask, nil, http.StatusOK)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+suite.mockTaskID.Hex(), nil)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Params = gin.Params{{Key: "id", Value: suite.mockTaskID.Hex()}}
	ctx.Request = req

	ctx.Set("user_id", suite.userID)
	ctx.Set("is_admin", suite.is_admin)

	suite.controller.GetTasksById(ctx)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestUpdateTasksById() {
	// Arrange

	mockTask := Domain.Task{
		Title:       "Updated Task",
		Description: "Updated description",
		DueDate:     time.Now(),
		Status:      "Completed",
		User_ID:     suite.userID,
	}

	suite.mockTaskUseCase.On("UpdateTasksById", mock.Anything, suite.mockTaskID, mockTask, mock.Anything).Return(mockTask, nil, http.StatusOK)

	// Act
	reqBody, _ := json.Marshal(mockTask)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+suite.mockTaskID.Hex(), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Params = gin.Params{{Key: "id", Value: suite.mockTaskID.Hex()}}
	ctx.Request = req

	ctx.Set("user_id", suite.userID)
	ctx.Set("is_admin", suite.is_admin)
	suite.controller.UpdateTasksById(ctx)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

func (suite *TaskControllerTestSuite) TestDeleteTasksById() {
	// Arrange

	suite.mockTaskUseCase.On("DeleteTasksById", mock.Anything, suite.mockTaskID, mock.Anything).Return(nil, http.StatusOK)

	// Act
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+suite.mockTaskID.Hex(), nil)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Params = gin.Params{{Key: "id", Value: suite.mockTaskID.Hex()}}
	ctx.Request = req

	ctx.Set("user_id", suite.userID)
	suite.controller.DeleteTasksById(ctx)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// TestTaskController runs the task controller test suite
func TestTaskController(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}
