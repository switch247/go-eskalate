package data

import (

	// "strconv"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "time"

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
	// tasks      map[string]*models.Task
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
	// collection.Drop(context.TODO()) //uncomment this tho drop collection
	ts := &TaskService{
		client:     client,
		validator:  validator.New(),
		DataBase:   DataBase,
		collection: collection,
	}

	return ts, nil
}

// create task
func (ts *TaskService) CreateTasks(task *models.Task) (models.Task, error, int) {
	// this needs some rework
	insertResult, err := ts.collection.InsertOne(context.TODO(), task)
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err, 500
	}

	fetched, err, _ := ts.GetTasksById(insertResult.InsertedID.(primitive.ObjectID))
	if err != nil {
		// this should have status 500
		return models.Task{}, errors.New("Task Not Created"), 500
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return fetched, nil, 201
}

// get all tasks
func (ts *TaskService) GetTasks() ([]*models.Task, error, int) {
	// ts.mu.RLock()
	// defer ts.mu.RUnlock()
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return nil, err1, 500
	}

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

// get task by id
func (ts *TaskService) GetTasksById(id primitive.ObjectID) (models.Task, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return models.Task{}, err1, 500
	}

	filter := bson.D{{"_id", id}}
	var result models.Task
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.Task{}, errors.New("Task not found"), http.StatusNotFound
	}
	return result, nil, 200
}

// update task by id
// used transactions for this one
func (ts *TaskService) UpdateTasksById(id primitive.ObjectID, task models.Task) (models.Task, error, int) {
	// Start a session
	session, err := ts.client.StartSession()
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err, 500
	}
	defer session.EndSession(context.Background())

	var NewTask models.Task
	statusCode := 200

	// Execute the transaction
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Start transaction
		err = session.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Retrieve the existing task
		NewTask, err, statusCode = ts.GetTasksById(id)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			return errors.New("task does not exist")
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

		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$set", bson.D{
				{"title", NewTask.Title},
				{"description", NewTask.Description},
				{"status", NewTask.Status},
				{"due_date", NewTask.DueDate},
			}},
		}

		updateResult, err := ts.collection.UpdateOne(sc, filter, update)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			fmt.Println(err)
			return err
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		// Commit the transaction
		err = session.CommitTransaction(sc)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		return models.Task{}, err, 500
	}

	return NewTask, nil, statusCode
}

// delete task by id
func (ts *TaskService) DeleteTasksById(id primitive.ObjectID) (error, int) {
	filter := bson.D{{"_id", id}}
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
