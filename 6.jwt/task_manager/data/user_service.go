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

type UserService struct {
	validator  *validator.Validate
	client     *mongo.Client
	DataBase   *mongo.Database
	collection *mongo.Collection
}

type userService interface {
	GetUsersByEmail(email string) (models.User, error, int)
	GetUsers() ([]models.User, error)
	CreateUsers(user *models.User) (models.User, error)
	GetUsersById(id string) (models.User, error)
	UpdateUsersById(id string, user models.User) (models.User, error)
	DeleteUsersById(id string) error
}

func NewUserService() (*UserService, error) {
	client, err := config.GetClient()
	DataBase := client.Database("test")
	_collection := DataBase.Collection("users")
	if err == nil {
		return &UserService{
			validator:  validator.New(),
			client:     client,
			DataBase:   DataBase,
			collection: _collection,
		}, nil
	} else {
		return nil, err
	}

}

// create user
func (ts *UserService) CreateUsers(user *models.User) (models.User, error, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := ts.client.StartSession()
	if err != nil {
		return models.User{}, err, 500
	}
	defer session.EndSession(ctx)

	resultTask := models.User{}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		insertResult, err := ts.collection.InsertOne(sc, user)
		if err != nil {
			fmt.Println(err)
			session.AbortTransaction(sc)
			return err
		}
		// Fetch the inserted user
		var fetched models.User
		err = ts.collection.FindOne(sc, bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)

		// fetched, err, _ := ts.GetUsersById(insertResult.InsertedID.(primitive.ObjectID))
		if err != nil {
			fmt.Println(err)
			session.AbortTransaction(sc)
			return errors.New("User Not Created")
		}

		if fetched.Email != user.Email {
			session.AbortTransaction(sc)
			return errors.New("User Not Created")
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
		return models.User{}, err, 500
	}

	fmt.Println("Inserted a single document: ", resultTask.ID)
	return resultTask, nil, 201
}

// get all users
func (ts *UserService) GetUsers() ([]*models.OmitedUser, error, int) {
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
	var results []*models.OmitedUser

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := ts.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
		return []*models.OmitedUser{}, err, 0
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem models.OmitedUser
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err.Error())
			// #handelthislater
			// should this say there was a decoding error and return?
			return []*models.OmitedUser{}, err, 500
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return []*models.OmitedUser{}, err, 500
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil, 200
}

// get user by id
func (ts *UserService) GetUsersById(id primitive.ObjectID) (models.User, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return models.User{}, err1, 500
	}

	filter := bson.D{{"_id", id}}
	var result models.User
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.User{}, errors.New("User not found"), http.StatusNotFound
	}
	return result, nil, 200
}

// get user by id
func (ts *UserService) GetUsersByEmail(email string) (models.User, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"email", 1}},
	})
	if err1 != nil {
		return models.User{}, err1, 500
	}

	filter := bson.D{{"email", email}}
	var result models.User
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.User{}, errors.New("User not found"), http.StatusNotFound
	}
	return result, nil, 200
}

// update user by id
// used transactions for this one
// func (ts *UserService) UpdateUsersById(id primitive.ObjectID, user models.User) (models.User, error, int) {
// 	// Start a session
// 	session, err := ts.client.StartSession()
// 	if err != nil {
// 		fmt.Println(err)
// 		return models.User{}, err, 500
// 	}
// 	defer session.EndSession(context.Background())

// 	var NewUser models.User
// 	statusCode := 200

// 	// Execute the transaction
// 	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
// 		// Start transaction
// 		err = session.StartTransaction()
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}

// 		// Retrieve the existing user
// 		NewUser, err, statusCode = ts.GetUsersById(id)
// 		if err != nil {
// 			_ = session.AbortTransaction(sc) // Roll back the transaction on error
// 			return errors.New("user does not exist")
// 		}

// 		// Update only the specified fields
// 		if user.Email != "" {
// 			NewUser.Email = user.Email
// 		}

// 		filter := bson.D{{"_id", id}}
// 		update := bson.D{
// 			{"$set", bson.D{
// 				{"title", NewUser.Email},
// 			}},
// 		}

// 		updateResult, err := ts.collection.UpdateOne(sc, filter, update)
// 		if err != nil {
// 			_ = session.AbortTransaction(sc) // Roll back the transaction on error
// 			fmt.Println(err)
// 			return err
// 		}
// 		if updateResult.ModifiedCount == 0 {
// 			return errors.New("user does not exist")
// 		}

// 		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

// 		// Commit the transaction
// 		err = session.CommitTransaction(sc)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return models.User{}, err, 500
// 	}

// 	return NewUser, nil, statusCode
// }

// // delete user by id
// func (ts *UserService) DeleteUsersById(id primitive.ObjectID) (error, int) {
// 	filter := bson.D{{"_id", id}}

// 	deleteResult, err := ts.collection.DeleteOne(context.TODO(), filter)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err, 500
// 	}
// 	if deleteResult.DeletedCount != 0 {
// 		return errors.New("User does not exist"), http.StatusNotFound
// 	}
// 	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
// 	return nil, 200

// }
