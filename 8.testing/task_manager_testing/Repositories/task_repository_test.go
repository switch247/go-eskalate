package Repositories_test

import (
	"context"
	"main/mongo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/mongo"

	"main/Domain"
	repo "main/Repositories"
	"main/utils"
)

type TaskRepositoryTestSuite struct {
	suite.Suite
	client     mongo.Client
	database   mongo.Database
	collection mongo.Collection
	repo       Domain.TaskRepository
}

func (suite *TaskRepositoryTestSuite) SetupTest() {
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
	collection := database.Collection("tasks")

	repo, err := repo.NewTaskRepository(client, database, collection)
	if err != nil {
		suite.T().Errorf("failed to create authRepository: %v", err)
		return
	}

	suite.client = client
	suite.database = database
	suite.collection = collection
	suite.repo = repo
}

func (suite *TaskRepositoryTestSuite) TearDownTest() {
	// Clean up the database after the tests
	// if err := suite.database.Drop(context.Background()); err != nil {
	// 	suite.T().Errorf("failed to drop database: %v", err)
	// }
}

func (suite *TaskRepositoryTestSuite) TestCreateTask() {
	// Create a new task
	task := &Domain.Task{
		Title:       "Test Task",
		Description: "This is a test task.",
		Status:      "To Do",
		DueDate:     time.Now().Add(time.Hour * 24 * 7),
		User_ID:     primitive.NewObjectID().Hex(),
	}

	createdTask, err, statusCode := suite.repo.CreateTasks(context.Background(), task)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, statusCode)
	assert.Equal(suite.T(), task.Title, createdTask.Title)
	assert.Equal(suite.T(), task.Description, createdTask.Description)
	assert.Equal(suite.T(), task.Status, createdTask.Status)
	assert.Equal(suite.T(), task.User_ID, createdTask.User_ID)
}

func (suite *TaskRepositoryTestSuite) TestGetTasks() {
	// Create a new task
	obj_id := primitive.NewObjectID()
	task := &Domain.Task{
		Title:       "Test Task 1",
		Description: "This is a test task.",
		Status:      "To Do",
		DueDate:     time.Now().Add(time.Hour * 24 * 7),
		User_ID:     utils.ObjectIdToString(obj_id),
	}

	_, err, _ := suite.repo.CreateTasks(context.Background(), task)
	assert.NoError(suite.T(), err)
	// flush, err := suite.collection.DeleteMany(context.Background(), primitive.M{})

	// Get all tasks
	user := Domain.OmitedUser{
		ID:       obj_id,
		Is_Admin: true,
	}

	tasks, err, statusCode := suite.repo.GetTasks(context.Background(), user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Greater(suite.T(), len(tasks), 0)
}

func (suite *TaskRepositoryTestSuite) TestGetTaskById() {
	// Create a new task
	task := &Domain.Task{
		Title:       "Test Task 2",
		Description: "This is a test task.",
		Status:      "To Do",
		DueDate:     time.Now().Add(time.Hour * 24 * 7),
		User_ID:     primitive.NewObjectID().Hex(),
	}

	createdTask, err, _ := suite.repo.CreateTasks(context.Background(), task)
	assert.NoError(suite.T(), err)

	// Get the task by ID
	obj_id, _ := utils.StringToObjectId(createdTask.User_ID)
	user := Domain.OmitedUser{
		ID:       obj_id,
		Is_Admin: true,
	}

	fetchedTask, err, statusCode := suite.repo.GetTasksById(context.Background(), createdTask.ID, user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
	assert.Equal(suite.T(), createdTask.Title, fetchedTask.Title)
	assert.Equal(suite.T(), createdTask.Description, fetchedTask.Description)
	assert.Equal(suite.T(), createdTask.Status, fetchedTask.Status)
	assert.Equal(suite.T(), createdTask.User_ID, fetchedTask.User_ID)
}

func (suite *TaskRepositoryTestSuite) TestUpdateTask() {
	// Create a new task
	task := &Domain.Task{
		Title:       "Test Task 3",
		Description: "This is a test task.",
		Status:      "To Do",
		DueDate:     time.Now().Add(time.Hour * 24 * 7),
		User_ID:     primitive.NewObjectID().Hex(),
	}

	createdTask, err, _ := suite.repo.CreateTasks(context.Background(), task)
	assert.NoError(suite.T(), err)

	// Update the task
	updatedTask := Domain.Task{
		Title:       "Updated Test Task",
		Description: "This is an updated test task.",
		Status:      "In Progress",
		DueDate:     time.Now().Add(time.Hour * 24 * 14),
	}
	obj_id, _ := utils.StringToObjectId(createdTask.User_ID)
	user := Domain.OmitedUser{
		ID:       obj_id,
		Is_Admin: true,
	}

	_, err, statusCode := suite.repo.UpdateTasksById(context.Background(), createdTask.ID, updatedTask, user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
}

func (suite *TaskRepositoryTestSuite) TestDeleteTask() {
	// Create a new task
	task := &Domain.Task{
		Title:       "Test Task 4",
		Description: "This is a test task.",
		Status:      "To Do",
		DueDate:     time.Now().Add(time.Hour * 24 * 7),
		User_ID:     primitive.NewObjectID().Hex(),
	}

	createdTask, err, _ := suite.repo.CreateTasks(context.Background(), task)
	assert.NoError(suite.T(), err)

	// Delete the task
	obj_id, _ := utils.StringToObjectId(createdTask.User_ID)
	user := Domain.OmitedUser{
		ID:       obj_id,
		Is_Admin: true,
	}

	err, statusCode := suite.repo.DeleteTasksById(context.Background(), createdTask.ID, user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, statusCode)
}

func TestTaskRepository(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
