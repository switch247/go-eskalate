package Repositories

import (

	// "strconv"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	// "os"

	"main/Domain"
	"main/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator"
)

type taskRepository struct {
	validator  *validator.Validate
	client     *mongo.Client
	DataBase   *mongo.Database
	collection *mongo.Collection
}

func NewTaskRepository(client *mongo.Client, DataBase *mongo.Database, _collection *mongo.Collection) (*taskRepository, error) {

	// collection.Drop(context.TODO()) //uncomment this tho drop collection
	ts := &taskRepository{
		client:     client,
		validator:  validator.New(),
		DataBase:   DataBase,
		collection: _collection,
	}

	return ts, nil

}

// create task
func (ts *taskRepository) CreateTasks(ctx context.Context, task *Domain.Task) (Domain.Task, error, int) {

	session, err := ts.client.StartSession()
	if err != nil {
		return Domain.Task{}, err, 500
	}
	defer session.EndSession(ctx)

	resultTask := Domain.Task{}
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
		var fetched Domain.Task
		err = ts.collection.FindOne(sc, bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)

		// fetched, err, _ := ts.GetTasksById(insertResult.InsertedID.(primitive.ObjectID))
		if err != nil {
			fmt.Println(err)
			session.AbortTransaction(sc)
			return errors.New("Task Not Created")
		}
		// var date_not_match = !fetched.DueDate.Equal(task.DueDate.In(fetched.DueDate.Location()))
		if fetched.Title != task.Title || fetched.Description != task.Description || fetched.Status != task.Status {
			session.AbortTransaction(sc)
			return errors.New("Task Not Created Properly")
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
		return Domain.Task{}, err, 500
	}

	fmt.Println("Inserted a single document: ", resultTask.ID)
	return resultTask, nil, 201
}

// get all tasks
func (ts *taskRepository) GetTasks(ctx context.Context, user Domain.OmitedUser) ([]*Domain.Task, error, int) {
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
	var filter bson.D
	if user.Is_Admin == true {
		filter = bson.D{{}}
	} else {
		userId := primitive.ObjectID.Hex(user.ID)

		filter = bson.D{{"user_id", userId}}
	}

	// Here's an array in which you can store the decoded documents
	var results []*Domain.Task

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := ts.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
		return []*Domain.Task{}, err, 0
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Domain.Task
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err.Error())
			// #handelthislater
			// should this say there was a decoding error and return?
			return []*Domain.Task{}, err, 500
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return []*Domain.Task{}, err, 500
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil, 200
}

// get task by id
func (ts *taskRepository) GetTasksById(ctx context.Context, id primitive.ObjectID, user Domain.OmitedUser) (Domain.Task, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return Domain.Task{}, err1, 500
	}

	filter := bson.D{{"_id", id}}
	var result Domain.Task
	err := ts.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Domain.Task{}, errors.New("Task not found"), http.StatusNotFound
	}

	val, err := primitive.ObjectIDFromHex(result.User_ID)
	if err != nil {
		return Domain.Task{}, errors.New("task id invalid"), 404
	}
	if user.Is_Admin == false && val != user.ID {
		return Domain.Task{}, errors.New("permission denied"), http.StatusUnauthorized

	}
	return result, nil, 200
}

// update task by id
// used transactions for this one
func (ts *taskRepository) UpdateTasksById(ctx context.Context, id primitive.ObjectID, task Domain.Task, user Domain.OmitedUser) (Domain.Task, error, int) {
	// Start a session
	session, err := ts.client.StartSession()
	if err != nil {
		fmt.Println(err)
		return Domain.Task{}, err, 500
	}
	defer session.EndSession(ctx)

	var NewTask Domain.Task
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
		NewTask, err, statusCode = ts.GetTasksById(ctx, id, user)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			return err
		}
		val, err := primitive.ObjectIDFromHex(NewTask.User_ID)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			return errors.New("task id invalid")
		}

		if user.Is_Admin == false && user.ID != val {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			statusCode = http.StatusUnauthorized
			return errors.New("user does not have permission to update this task")
		}
		// Update only the specified fields
		utils.UpdateFields(task, &NewTask)

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
		return Domain.Task{}, err, 500
	}

	return NewTask, nil, statusCode
}

// delete task by id
func (ts *taskRepository) DeleteTasksById(ctx context.Context, id primitive.ObjectID, user Domain.OmitedUser) (error, int) {
	filter := bson.D{{"_id", id}}

	// Retrieve the existing task
	NewTask, err, statusCode := ts.GetTasksById(ctx, id, user)
	if err != nil {
		return err, statusCode
	}

	val, err := primitive.ObjectIDFromHex(NewTask.User_ID)
	if err != nil {
		return errors.New("task id invalid"), 404
	}

	if user.Is_Admin == false && user.ID != val {
		return errors.New("user does not have permission to update this task"), 403
	}

	deleteResult, err := ts.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("Task does not exist"), http.StatusNotFound
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil, 200

}
