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
)

type UserUseCaseTestSuite struct {
	suite.Suite
	mockUserRepo *mocks.UserRepository
	userUseCase  Domain.UserUseCases
	ctx          *gin.Context
}

func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.UserRepository)
	var err error
	suite.ctx = &gin.Context{}
	suite.userUseCase, err = UseCases.NewUserUseCase(suite.mockUserRepo)
	assert.NoError(suite.T(), err)
}

func (suite *UserUseCaseTestSuite) TestGetUsers() {
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

	suite.mockUserRepo.On("GetUsers", mock.Anything).Return(users, nil, 200).Once()

	// Act
	result, err, statusCode := suite.userUseCase.GetUsers(suite.ctx)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), users, result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserUseCaseTestSuite) TestGetUsersById() {
	// Arrange
	user := Domain.OmitedUser{
		ID:    primitive.NewObjectID(),
		Email: "test@example.com",
	}

	suite.mockUserRepo.On("GetUsersById", mock.Anything, user.ID, user).Return(user, nil, 200).Once()

	// Act
	result, err, statusCode := suite.userUseCase.GetUsersById(suite.ctx, user.ID, user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), user, result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserUseCaseTestSuite) TestCreateUsers() {
	// Arrange
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Is_Admin: true,
	}

	createdUser := Domain.OmitedUser{
		ID:    primitive.NewObjectID(),
		Email: "test@example.com",
	}

	suite.mockUserRepo.On("CreateUsers", mock.Anything, user).Return(createdUser, nil, 201).Once()

	// Act
	result, err, statusCode := suite.userUseCase.CreateUsers(suite.ctx, user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, statusCode)
	assert.Equal(suite.T(), createdUser, result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserUseCaseTestSuite) TestUpdateUsersById() {
	// Arrange
	userID := primitive.NewObjectID()
	updatedUser := Domain.User{
		Email: "updated@example.com",
	}

	currentUser := Domain.OmitedUser{
		ID:    userID,
		Email: "test@example.com",
	}

	updatedOmitedUser := Domain.OmitedUser{
		ID:    userID,
		Email: "updated@example.com",
	}

	suite.mockUserRepo.On("UpdateUsersById", mock.Anything, userID, updatedUser, currentUser).Return(updatedOmitedUser, nil, 200).Once()

	// Act
	result, err, statusCode := suite.userUseCase.UpdateUsersById(suite.ctx, userID, updatedUser, currentUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), updatedOmitedUser, result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserUseCaseTestSuite) TestDeleteUsersById() {
	// Arrange
	userID := primitive.NewObjectID()
	currentUser := Domain.OmitedUser{
		ID:    userID,
		Email: "test@example.com",
	}

	suite.mockUserRepo.On("DeleteUsersById", mock.Anything, userID, currentUser).Return(nil, 200).Once()

	// Act
	err, statusCode := suite.userUseCase.DeleteUsersById(suite.ctx, userID, currentUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserUseCaseTestSuite) TestGetUsersError() {
	// Arrange
	suite.mockUserRepo.On("GetUsers", mock.Anything).Return(nil, errors.New("get users error"), 500).Once()

	// Act
	result, err, statusCode := suite.userUseCase.GetUsers(suite.ctx)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 500, statusCode)
	assert.Nil(suite.T(), result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func TestUserUseCase(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}
