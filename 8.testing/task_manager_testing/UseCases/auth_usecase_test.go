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

type AuthUseCaseTestSuite struct {
	suite.Suite
	ctx          *gin.Context
	mockAuthRepo *mocks.AuthRepository
	authUseCase  Domain.AuthUseCase
}

func (suite *AuthUseCaseTestSuite) SetupTest() {
	suite.mockAuthRepo = new(mocks.AuthRepository)
	var err error
	suite.ctx = &gin.Context{}
	suite.authUseCase, err = UseCases.NewAuthUseCase(suite.mockAuthRepo)
	assert.NoError(suite.T(), err)
}

func (suite *AuthUseCaseTestSuite) TestLogin() {
	// Arrange
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockAuthRepo.On("Login", mock.Anything, user).Return("token", nil, 200).Once()

	// Act
	token, err, statusCode := suite.authUseCase.Login(suite.ctx, user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), "token", token)
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestRegister() {
	// Arrange
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	createdUser := Domain.OmitedUser{
		ID:    primitive.NewObjectID(),
		Email: "test@example.com",
	}

	suite.mockAuthRepo.On("Register", mock.Anything, user).Return(createdUser, nil, 201).Once()

	// Act
	result, err, statusCode := suite.authUseCase.Register(suite.ctx, user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, statusCode)
	assert.Equal(suite.T(), createdUser, result)
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestLoginError() {
	// Arrange
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockAuthRepo.On("Login", mock.Anything, user).Return("", errors.New("login error"), 401).Once()

	// Act
	token, err, statusCode := suite.authUseCase.Login(suite.ctx, user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 401, statusCode)
	assert.Empty(suite.T(), token)
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestRegisterError() {
	// Arrange
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	suite.mockAuthRepo.On("Register", mock.Anything, user).Return(Domain.OmitedUser{}, errors.New("registration error"), 500).Once()

	// Act
	result, err, statusCode := suite.authUseCase.Register(suite.ctx, user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 500, statusCode)
	assert.Equal(suite.T(), Domain.OmitedUser{}, result)
	suite.mockAuthRepo.AssertExpectations(suite.T())
}

func TestAuthUseCase(t *testing.T) {
	// suite := &AuthUseCaseTestSuite{}
	// suite.SetupTest()
	suite.Run(t, new(AuthUseCaseTestSuite))
}
