package Repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/Domain"
	repo "main/Repositories"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	repo       Domain.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// Initialize the Mongo client and database
	url := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		suite.T().Errorf("failed to create Mongo client: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		suite.T().Errorf("Failed to connect to MongoDB: %v", err)

		return
	}

	database := client.Database("testing")
	collection := database.Collection("user")

	repo, err := repo.NewUserRepository(client, database, collection)
	if err != nil {
		suite.T().Errorf("failed to create authRepository: %v", err)
		return
	}

	suite.client = client
	suite.database = database
	suite.collection = collection
	suite.repo = repo
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	// Clean up the database after the tests
	if err := suite.database.Drop(context.Background()); err != nil {
		suite.T().Errorf("failed to drop database: %v", err)
	}
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	// Create a new user
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	createdUser, err, statusCode := suite.repo.CreateUsers(context.Background(), user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), user.Email, createdUser.Email)
	assert.NotEqual(suite.T(), user.Password, createdUser.Password)
}

func (suite *UserRepositoryTestSuite) TestGetUsers() {
	// Create some test users
	user1 := &Domain.User{
		Email:    "test1@example.com",
		Password: "password123",
	}
	user2 := &Domain.User{
		Email:    "test2@example.com",
		Password: "password456",
	}

	_, _, _ = suite.repo.CreateUsers(context.Background(), user1)
	_, _, _ = suite.repo.CreateUsers(context.Background(), user2)

	users, err, statusCode := suite.repo.GetUsers(context.Background())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), 2, len(users))
}

func (suite *UserRepositoryTestSuite) TestGetUserByID() {
	// Create a test user
	obj_id := primitive.NewObjectID()
	user := &Domain.User{
		ID:       obj_id,
		Email:    "test@example.com",
		Password: "password123",
		Is_Admin: true,
	}

	createdUser, _, _ := suite.repo.CreateUsers(context.Background(), user)

	fetchedUser, err, statusCode := suite.repo.GetUsersById(context.Background(), obj_id, createdUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), createdUser.Email, fetchedUser.Email)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	// Create a test user
	obj_id := primitive.NewObjectID()
	user := &Domain.User{
		ID:       obj_id,
		Email:    "test@example.com",
		Password: "password123",
		Is_Admin: true,
	}

	createdUser, _, _ := suite.repo.CreateUsers(context.Background(), user)

	// Update the user
	updatedUser := &Domain.User{
		Email: "updated@example.com",
	}

	updatedUserResult, err, statusCode := suite.repo.UpdateUsersById(context.Background(), createdUser.ID, *updatedUser, createdUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), updatedUser.Email, updatedUserResult.Email)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser() {
	// Create a test user
	user := &Domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	createdUser, _, _ := suite.repo.CreateUsers(context.Background(), user)

	err, statusCode := suite.repo.DeleteUsersById(context.Background(), createdUser.ID, createdUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)

	_, err, statusCode = suite.repo.GetUsersById(context.Background(), createdUser.ID, createdUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 404, statusCode)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
