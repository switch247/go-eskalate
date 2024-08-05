package data

import (

	// "strconv"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	// "os"
	"time"

	"main/config"
	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
}

func NewTaskService() (*TaskService, error) {

	// Set client options
	// Connect to MongoDB
	// Check the connection
	var collection *mongo.Collection
	var DataBase *mongo.Database
	client, err := config.GetClient()
	if err == nil {
		DataBase = client.Database("test")
		collection = DataBase.Collection("tasks")
		// collection.Drop(context.TODO()) //uncomment this tho drop collection
		ts := &TaskService{
			client:     client,
			validator:  validator.New(),
			DataBase:   DataBase,
			collection: collection,
		}

		return ts, nil
	} else {
		return nil, err
	}
}

// create task
func (ts *TaskService) CreateTasks(task *models.Task) (models.Task, error, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := ts.client.StartSession()
	if err != nil {
		return models.Task{}, err, 500
	}
	defer session.EndSession(ctx)

	resultTask := models.Task{}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		insertResult, err := ts.collection.InsertOne(sc, task)
		if err != nil {
			fmt.Println(err)
			session.AbortTransaction(sc)
			return err
		}
		// Fetch the inserted task
		var fetched models.Task
		err = ts.collection.FindOne(sc, bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)

		// fetched, err, _ := ts.GetTasksById(insertResult.InsertedID.(primitive.ObjectID))
		if err != nil {
			fmt.Println(err)
			session.AbortTransaction(sc)
			return errors.New("Task Not Created")
		}
		var date_not_match = !fetched.DueDate.Equal(task.DueDate.In(fetched.DueDate.Location()))
		if fetched.Title != task.Title || fetched.Description != task.Description || fetched.Status != task.Status || date_not_match {
			session.AbortTransaction(sc)
			return errors.New("Task Not Created")
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			fmt.Println(err)

			return err
		}

		resultTask = fetched
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return models.Task{}, err, 500
	}

	fmt.Println("Inserted a single document: ", resultTask.ID)
	return resultTask, nil, 201
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
func (ts *TaskService) UpdateTasksById(id primitive.ObjectID, task models.Task, userid primitive.ObjectID) (models.Task, error, int) {
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
		val, err := primitive.ObjectIDFromHex(NewTask.User_ID)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			return errors.New("task does not exist")
		}

		if userid != val {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			statusCode = http.StatusUnauthorized
			return errors.New("user does not have permission to update this task")
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
		if task.User_ID != "" {
			fmt.Println("user id", task.User_ID)
			NewTask.User_ID = task.User_ID
		}

		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$set", bson.D{
				{"title", NewTask.Title},
				{"description", NewTask.Description},
				{"status", NewTask.Status},
				{"due_date", NewTask.DueDate},
				{"user_id", NewTask.User_ID},
			}},
		}

		updateResult, err := ts.collection.UpdateOne(sc, filter, update)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			fmt.Println(err)
			return err
		}

		if updateResult.ModifiedCount == 0 {
			fmt.Println("Task does not exist")
			return errors.New("task does not exist")
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
func (ts *TaskService) DeleteTasksById(id primitive.ObjectID, userid primitive.ObjectID) (error, int) {
	filter := bson.D{{"_id", id}}

	// Retrieve the existing task
	NewTask, err, statusCode := ts.GetTasksById(id)
	if err != nil {
		return errors.New("task does not exist"), statusCode
	}

	val, err := primitive.ObjectIDFromHex(NewTask.User_ID)
	if err != nil {
		return errors.New("task does not exist"), 404
	}

	if userid != val {
		return errors.New("user does not have permission to update this task"), 403
	}

	deleteResult, err := ts.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	if deleteResult.DeletedCount != 0 {
		return errors.New("Task does not exist"), http.StatusNotFound
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil, 200

}
