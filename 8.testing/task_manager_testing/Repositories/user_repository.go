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
	"main/Infrastructure"
	"main/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator"
)

type userRepository struct {
	validator  *validator.Validate
	client     *mongo.Client
	DataBase   *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, DataBase *mongo.Database, _collection *mongo.Collection) (*userRepository, error) {

	return &userRepository{
		validator:  validator.New(),
		client:     client,
		DataBase:   DataBase,
		collection: _collection,
	}, nil

}

// create user
func (as *userRepository) CreateUsers(ctx context.Context, user *Domain.User) (Domain.OmitedUser, error, int) {

	// Check if user email is taken
	existingUserFilter := bson.D{{"email", user.Email}}
	existingUserCount, err := as.collection.CountDocuments(ctx, existingUserFilter)
	if err != nil {
		return Domain.OmitedUser{}, err, 500
	}
	if existingUserCount > 0 {
		return Domain.OmitedUser{}, errors.New("Email is already taken"), http.StatusBadRequest
	}
	// User registration logic
	hashedPassword, err := Infrastructure.GenerateFromPasswordCustom(user.Password)
	if err != nil {
		return Domain.OmitedUser{}, err, 500
	}
	user.Password = string(hashedPassword)
	insertResult, err := as.collection.InsertOne(ctx, user)
	if err != nil {
		return Domain.OmitedUser{}, err, 500
	}
	// Fetch the inserted task
	var fetched Domain.OmitedUser
	err = as.collection.FindOne(context.TODO(), bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)
	if err != nil {
		fmt.Println(err)
		return Domain.OmitedUser{}, errors.New("User Not Created"), 500
	}
	if fetched.Email != user.Email {
		return Domain.OmitedUser{}, errors.New("User Not Created"), 500
	}
	fetched.Password = ""
	return fetched, nil, 200
}

// get all users
func (ts *userRepository) GetUsers(ctx context.Context) ([]*Domain.OmitedUser, error, int) {
	// ts.mu.RLock()
	// defer ts.mu.RUnlock()
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
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
	var results []*Domain.OmitedUser

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := ts.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
		return []*Domain.OmitedUser{}, err, 0
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var elem Domain.OmitedUser
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err.Error())
			// #handelthislater
			// should this say there was a decoding error and return?
			return []*Domain.OmitedUser{}, err, 500
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return []*Domain.OmitedUser{}, err, 500
	}

	// Close the cursor once finished
	cur.Close(ctx)

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil, 200
}

// get user by id
func (ts *userRepository) GetUsersById(ctx context.Context, id primitive.ObjectID, user Domain.OmitedUser) (Domain.OmitedUser, error, int) {
	// Create an index on the "_id" field
	_, err1 := ts.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return Domain.OmitedUser{}, err1, 500
	}
	var filter bson.D
	if user.Is_Admin == false {
		fmt.Println("user is not admin")
		userIdString := utils.ObjectIdToString(user.ID)
		filter = bson.D{{"_id", id}, {"user_id", userIdString}}

	} else {
		filter = bson.D{{"_id", id}}
	}
	var result Domain.OmitedUser
	err := ts.collection.FindOne(ctx, filter).Decode(&result)
	// # handel this later
	if err != nil {
		return Domain.OmitedUser{}, errors.New("User not found"), http.StatusNotFound
	}
	if user.Is_Admin == false && result.ID != user.ID {
		return Domain.OmitedUser{}, errors.New("permission denied"), http.StatusForbidden

	}
	return result, nil, 200
}

// update user by id
// used transactions for this one

func (ts *userRepository) UpdateUsersById(ctx context.Context, id primitive.ObjectID, user Domain.User, curentuser Domain.OmitedUser) (Domain.OmitedUser, error, int) {
	// Start a session
	session, err := ts.client.StartSession()
	if err != nil {
		fmt.Println(err)
		return Domain.OmitedUser{}, err, 500
	}
	defer session.EndSession(ctx)

	var NewUser Domain.OmitedUser
	statusCode := 200

	// Execute the transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// Start transaction
		err = session.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Retrieve the existing user
		NewUser, err, statusCode = ts.GetUsersById(ctx, id, curentuser)
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
		return Domain.OmitedUser{}, err, 500
	}

	return NewUser, nil, statusCode
}

// delete user by id
func (ts *userRepository) DeleteUsersById(ctx context.Context, id primitive.ObjectID, user Domain.OmitedUser) (error, int) {

	filter := bson.D{{"_id", id}}
	if user.Is_Admin == false && user.ID != id {
		return errors.New("permision denied"), http.StatusForbidden
	}

	deleteResult, err := ts.collection.DeleteOne(ctx, filter)
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

// get user by email
// func (ts *userRepository) GetUsersByEmail(email string) (Domain.OmitedUser, error, int) {
// 	// Create an index on the "_id" field
// 	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
// 		Keys: bson.D{{"email", 1}},
// 	})
// 	if err1 != nil {
// 		return Domain.OmitedUser{}, err1, 500
// 	}

// 	filter := bson.D{{"email", email}}
// 	var result Domain.OmitedUser
// 	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
// 	if err != nil {
// 		return Domain.OmitedUser{}, errors.New("User not found"), http.StatusNotFound
// 	}
// 	return result, nil, 200
// }
