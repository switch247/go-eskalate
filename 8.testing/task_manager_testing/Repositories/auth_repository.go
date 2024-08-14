package Repositories

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"main/Domain"
	"main/Infrastructure"
	mongo "main/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/mongo"

	"github.com/go-playground/validator"
)

type authRepository struct {
	validator  *validator.Validate
	client     mongo.Client
	DataBase   mongo.Database
	collection mongo.Collection
}

func NewAuthRepository(client mongo.Client, DataBase mongo.Database, _collection mongo.Collection) (*authRepository, error) {
	// _collection.Drop(context.TODO()) //un comment this tho drop collection
	return &authRepository{
		validator:  validator.New(),
		client:     client,
		DataBase:   DataBase,
		collection: _collection,
	}, nil

}

func (au *authRepository) Login(ctx context.Context, user *Domain.User) (string, error, int) {
	// TODO: Implement user login logic
	filter := bson.D{{"email", user.Email}}
	var existingUser Domain.User
	err := au.collection.FindOne(ctx, filter).Decode(&existingUser)

	if err != nil || !Infrastructure.CompareHashAndPasswordCustom(existingUser.Password, user.Password) {
		fmt.Printf("Login Called:%v, %v", existingUser.Password, user.Password)
		hashedPassword, _ := Infrastructure.GenerateFromPasswordCustom(user.Password)
		fmt.Print(existingUser.Password == hashedPassword)
		return "", errors.New("Invalid credentials"), http.StatusBadRequest
	}

	// Generate JWT
	jwtToken, err := Infrastructure.SignJwt(existingUser)
	if err != nil {
		return "", err, 500
	}

	return jwtToken, nil, 200
}

func (au *authRepository) Register(ctx context.Context, user *Domain.User) (Domain.OmitedUser, error, int) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// Check if user email is taken
	existingUserFilter := bson.D{{"email", user.Email}}
	existingUserCount, err := au.collection.CountDocuments(ctx, existingUserFilter)
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
	InsertedID, err := au.collection.InsertOne(ctx, user)
	if err != nil {
		return Domain.OmitedUser{}, err, 500
	}
	// Fetch the inserted task
	var fetched Domain.OmitedUser
	err = au.collection.FindOne(context.TODO(), bson.D{{"_id", InsertedID.(primitive.ObjectID)}}).Decode(&fetched)
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
