package data

import (

	// "strconv"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"

	"errors"
	"main/models"
	"net/http"

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
	validator  *validator.Validate
	client     *mongo.Client
	DataBase   *mongo.Database
	collection *mongo.Collection
	tasks      map[string]*models.Task
	// mu        sync.RWMutex //i will add this back once i understand routines properly
}

func NewTaskService() (*TaskService, error) {
	// Read MONGO_CONNECTION_STRING from environment
	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING")
	// Set client options
	clientOptions := options.Client().ApplyURI(MONGO_CONNECTION_STRING)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	DataBase := client.Database("test")
	collection := DataBase.Collection("tasks")

	ts := &TaskService{
		client:     client,
		validator:  validator.New(),
		DataBase:   DataBase,
		collection: collection,
		tasks:      make(map[string]*models.Task),
	}
	ts.tasks["1"] = &models.Task{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"}
	ts.tasks["2"] = &models.Task{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"}
	ts.tasks["3"] = &models.Task{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"}
	return ts, nil
}

// get all tasks
func (ts *TaskService) GetTasks() ([]*models.Task, error, int) {
	// ts.mu.RLock()
	// defer ts.mu.RUnlock()

	// Pass these options to the Find method
	findOptions := options.Find()
	// findOptions.SetLimit(2)
	filter := bson.D{{}}

	// Here's an array in which you can store the decoded documents
	var results []*models.Task

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := ts.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
		return []*models.Task{}, err, 0
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.Task
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err.Error())
			// #handelthislater
			// should this say there was a decoding error and return?
			return []*models.Task{}, err, 500
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return []*models.Task{}, err, 500
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil, 200
}

// create task
func (ts *TaskService) CreateTasks(task *models.Task) (models.Task, error, int) {
	// this needs some rework
	insertResult, err := ts.collection.InsertOne(context.TODO(), task)
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err, 500
	}
	_, err, _ = ts.GetTasksById(task.ID)
	if err != nil {
		// this should have status 500
		return models.Task{}, errors.New("Task Not Created"), 500
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return *task, nil, 201
}

// get task by id
func (ts *TaskService) GetTasksById(id string) (models.Task, error, int) {
	filter := bson.D{{"id", id}}
	var result models.Task
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.Task{}, errors.New("Task not found"), http.StatusNotFound
	}
	return result, nil, 200
}

// update task by id
func (ts *TaskService) UpdateTasksById(id string, task models.Task) (models.Task, error, int) {
	NewTask, err, statusCode := ts.GetTasksById(id)

	if err != nil {
		return models.Task{}, errors.New("Task does not exists"), statusCode
	}
	// Update only the specified fields
	if task.Title != "" {
		NewTask.Title = task.Title
	}
	if task.Description != "" {
		NewTask.Description = task.Description
	}
	if task.Status != "" {
		NewTask.Status = task.Status
	}
	if !task.DueDate.IsZero() {
		NewTask.DueDate = task.DueDate
	}

	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"title", NewTask.Title},
			{"description", NewTask.Description},
			{"status", NewTask.Status},
			{"due_date", NewTask.DueDate},
		}},
	}
	updateResult, err := ts.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err, 500
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return NewTask, nil, 200

}

// delete task by id
func (ts *TaskService) DeleteTasksById(id string) (error, int) {
	filter := bson.D{{"id", id}}
	_, err1, status := ts.GetTasksById(id)
	if err1 != nil {
		fmt.Println(err1)
		return errors.New("Task does not exist"), status
	}
	deleteResult, err := ts.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil, 200

}
