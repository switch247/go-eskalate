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
	"main/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

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
func (as *UserService) CreateUsers(user *models.User) (models.OmitedUser, error, int) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// Check if user email is taken
	existingUserFilter := bson.D{{"email", user.Email}}
	existingUserCount, err := as.collection.CountDocuments(context.TODO(), existingUserFilter)
	if err != nil {
		return models.OmitedUser{}, err, 500
	}
	if existingUserCount > 0 {
		return models.OmitedUser{}, errors.New("Email is already taken"), http.StatusBadRequest
	}
	// User registration logic
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.OmitedUser{}, err, 500
	}
	user.Password = string(hashedPassword)
	insertResult, err := as.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.OmitedUser{}, err, 500
	}
	// Fetch the inserted task
	var fetched models.OmitedUser
	err = as.collection.FindOne(context.TODO(), bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)
	if err != nil {
		fmt.Println(err)
		return models.OmitedUser{}, errors.New("User Not Created"), 500
	}
	if fetched.Email != user.Email {
		return models.OmitedUser{}, errors.New("User Not Created"), 500
	}
	fetched.Password = ""
	return fetched, nil, 200
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
func (ts *UserService) GetUsersById(id primitive.ObjectID, user models.OmitedUser) (models.OmitedUser, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return models.OmitedUser{}, err1, 500
	}
	var filter bson.D
	if user.Is_Admin == false {
		userIdString := utils.ObjectIdToString(user.ID)
		filter = bson.D{{"_id", id}, {"user_id", userIdString}}

	} else {

		filter = bson.D{{"_id", id}}
	}
	var result models.OmitedUser
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.OmitedUser{}, errors.New("User not found"), http.StatusNotFound
	}
	if user.Is_Admin == false && result.ID != user.ID {
		return models.OmitedUser{}, errors.New("permission denied"), http.StatusForbidden

	}
	return result, nil, 200
}

// get user by id
func (ts *UserService) GetUsersByEmail(email string) (models.OmitedUser, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"email", 1}},
	})
	if err1 != nil {
		return models.OmitedUser{}, err1, 500
	}

	filter := bson.D{{"email", email}}
	var result models.OmitedUser
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.OmitedUser{}, errors.New("User not found"), http.StatusNotFound
	}
	return result, nil, 200
}

// update user by id
// used transactions for this one

func (ts *UserService) UpdateUsersById(id primitive.ObjectID, user models.User, curentuser models.OmitedUser) (models.OmitedUser, error, int) {
	// Start a session
	session, err := ts.client.StartSession()
	if err != nil {
		fmt.Println(err)
		return models.OmitedUser{}, err, 500
	}
	defer session.EndSession(context.Background())

	var NewUser models.OmitedUser
	statusCode := 200

	// Execute the transaction
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Start transaction
		err = session.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Retrieve the existing user
		NewUser, err, statusCode = ts.GetUsersById(id, curentuser)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			return err
		}

		// Update only the specified fields
		if user.Email != "" {
			NewUser.Email = user.Email
		}

		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$set", bson.D{
				{"title", NewUser.Email},
			}},
		}

		updateResult, err := ts.collection.UpdateOne(sc, filter, update)
		if err != nil {
			_ = session.AbortTransaction(sc) // Roll back the transaction on error
			fmt.Println(err)
			statusCode = 500
			return err
		}
		if updateResult.ModifiedCount == 0 {
			statusCode = 400
			return errors.New("user does not exist")
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		// Commit the transaction
		err = session.CommitTransaction(sc)
		if err != nil {
			statusCode = 500
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		return models.OmitedUser{}, err, 500
	}

	return NewUser, nil, statusCode
}

// delete user by id
func (ts *UserService) DeleteUsersById(id primitive.ObjectID, user models.OmitedUser) (error, int) {

	filter := bson.D{{"_id", id}}
	if user.Is_Admin == false && user.ID != id {
		return errors.New("permision denied"), http.StatusForbidden
	}

	deleteResult, err := ts.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("User does not exist"), http.StatusNotFound
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil, 200

}
