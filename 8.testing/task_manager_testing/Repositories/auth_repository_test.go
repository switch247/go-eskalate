package Repositories_test

import (
	"context"
	"testing"
	"time"

	"main/Domain"
	repo "main/Repositories"
	"main/mongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	// "go.mongodb.org/mongo-driver/mongo"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	client     mongo.Client
	database   mongo.Database
	collection mongo.Collection
	repo       Domain.AuthRepository
}

func (suite *AuthRepositoryTestSuite) SetupTest() {
	// Initialize the Mongo client and database
	url := "mongodb://localhost:27017"
	// clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.NewClient(url)
	if err != nil {
		suite.T().Errorf("failed to create Mongo client: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx)
	if err != nil {
		suite.T().Errorf("Failed to connect to MongoDB: %v", err)

		return
	}

	database := client.Database("testing")
	collection := database.Collection("user")

	repo, err := repo.NewAuthRepository(client, database, collection)
	if err != nil {
		suite.T().Errorf("failed to create authRepository: %v", err)
		return
	}

	suite.client = client
	suite.database = database
	suite.collection = collection
	suite.repo = repo
}

func (suite *AuthRepositoryTestSuite) TearDownTest() {
	// Clean up the database after the tests
	// if err := suite.database.Drop(context.Background()); err != nil {
	// 	suite.T().Errorf("failed to drop database: %v", err)
	// }
}

func (suite *AuthRepositoryTestSuite) TestRegister() {
	// Test the Register function with a valid user
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err, statusCode := suite.repo.Register(context.Background(), user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)

	// Test the case where the user's email is already taken
	user2 := &Domain.User{
		Email:    "test@example.com",
		Password: "password456",
	}

	_, err, statusCode = suite.repo.Register(context.Background(), user2)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 400, statusCode)
}

func (suite *AuthRepositoryTestSuite) TestLogin() {
	// Register a user
	new_user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err, statusCode := suite.repo.Register(context.Background(), new_user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	// Test the Login function with valid credentials
	jwtToken, err, statusCode := suite.repo.Login(context.Background(), new_user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)

	assert.NotEmpty(suite.T(), jwtToken)

	// Test the Login function with invalid credentials
	invalidUser := &Domain.User{
		Email:    "invalid@example.com",
		Password: "wrongpassword",
	}

	_, err, statusCode = suite.repo.Login(context.Background(), invalidUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 400, statusCode)

	// Test the Login function with wrong password
	wrongPasswordUser := &Domain.User{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err, statusCode = suite.repo.Login(context.Background(), wrongPasswordUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 400, statusCode)

}

func TestAuthRepository(t *testing.T) {
	_suite := new(AuthRepositoryTestSuite)
	suite.Run(t, _suite)
}
