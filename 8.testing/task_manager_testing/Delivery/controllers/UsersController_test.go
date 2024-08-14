package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/Delivery/controllers"
	"main/Domain"
	"main/mocks"
	"main/utils"
)

type UserControllerTestSuite struct {
	suite.Suite
	mockUserUseCase *mocks.UserUseCases
	userController  controllers.UserController
	userID          primitive.ObjectID
	is_admin        bool
}

func (suite *UserControllerTestSuite) SetupTest() {
	obj_id, _ := utils.StringToObjectId("66bc5f53f0e7bd3ca1d6bec9")
	suite.mockUserUseCase = new(mocks.UserUseCases)
	var err error
	suite.userController, err = controllers.NewUserController(suite.mockUserUseCase)
	assert.NoError(suite.T(), err)
	suite.userID = obj_id
	suite.is_admin = true
}

func (suite *UserControllerTestSuite) TestGetUsers() {
	// Arrange
	users := []*Domain.OmitedUser{
		{
			ID:    primitive.NewObjectID(),
			Email: "user1@example.com",
		},
		{
			ID:    primitive.NewObjectID(),
			Email: "user2@example.com",
		},
	}

	suite.mockUserUseCase.On("GetUsers", mock.Anything).Return(users, nil, 200).Once()

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	suite.userController.GetUsers(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	// assert.Equal(suite.T(), "{\"users\":[{\"ID\":\"", w.Body.String())
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetUsersError() {
	// Arrange
	suite.mockUserUseCase.On("GetUsers", mock.Anything).Return(nil, errors.New("get users error"), 500).Once()

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	user_id := primitive.NewObjectID().Hex()
	c.Set("user_id", user_id)
	// Act
	suite.userController.GetUsers(c)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	// assert.Equal(suite.T(), "{\"error\":\"get users error\"}", w.Body.String())
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetUser() {
	// Arrange
	user := Domain.OmitedUser{
		ID:    suite.userID,
		Email: "test@example.com",
	}

	suite.mockUserUseCase.On("GetUsersById", mock.Anything, user.ID, user).Return(user, nil, 200).Once()

	req, _ := http.NewRequest(http.MethodGet, "/users/"+user.ID.Hex(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: user.ID.Hex()},
	}
	c.Set("user_id", suite.userID)
	c.Set("is_admin", suite.is_admin)

	// Act
	suite.userController.GetUser(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	// assert.Equal(suite.T(), "{\"user\":{\"ID\":\""+user.ID.Hex()+"\",\"Name\":\"Test User\",\"Email\":\"test@example.com\"}}", w.Body.String())
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestCreateUser() {
	// Arrange
	user := Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	createdUser := Domain.OmitedUser{
		ID:    primitive.NewObjectID(),
		Email: "test@example.com",
	}

	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.mockUserUseCase.On("CreateUsers", c, &user).Return(createdUser, nil, 201).Once()

	// Act
	suite.userController.CreateUser(c)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	// assert.Equal(suite.T(), "{\"user\":{\"ID\":\"", w.Body.String())
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestUpdateUser() {
	// Arrange
	userID := primitive.NewObjectID()
	user := Domain.User{
		Email: "updated@example.com",
	}

	currentUser := Domain.OmitedUser{
		ID:       userID,
		Email:    "test@example.com",
		Is_Admin: true,
	}

	updatedUser := Domain.OmitedUser{
		ID:    userID,
		Email: "updated@example.com",
	}

	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPut, "/users/"+userID.Hex(), bytes.NewBuffer(jsonUser))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: userID.Hex()},
	}
	c.Set("user_id", suite.userID)
	c.Set("is_admin", suite.is_admin)

	suite.mockUserUseCase.On("UpdateUsersById", c, userID, user, currentUser).Return(updatedUser, nil, 200).Once()

	// Act
	suite.userController.UpdateUser(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	// assert.Equal(suite.T(), "{\"user\":{\"ID\":\""+userID.Hex()+"\",\"Name\":\"Updated User\",\"Email\":\"updated@example.com\"}}", w.Body.String())
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestDeleteUser() {
	// Arrange
	userID := primitive.NewObjectID()
	currentUser := Domain.OmitedUser{
		ID:       userID,
		Is_Admin: true,
		Email:    "test@example.com",
	}

	req, _ := http.NewRequest(http.MethodDelete, "/users/"+userID.Hex(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: userID.Hex()},
	}
	c.Set("user_id", suite.userID)
	c.Set("is_admin", suite.is_admin)

	suite.mockUserUseCase.On("DeleteUsersById", c, userID, currentUser).Return(nil, 200).Once()

	// Act
	suite.userController.DeleteUser(c)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	// assert.Equal(suite.T(), "{\"message\":\"User deleted successfully\"}", w.Body.String())
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
