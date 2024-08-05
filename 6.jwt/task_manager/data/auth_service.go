package data

import (
	// "strconv"
	"context"
	"errors"
	"fmt"
	"net/http"

	// "os"
	"time"

	"main/config"
	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
)

type AuthService struct {
	validator  *validator.Validate
	client     *mongo.Client
	DataBase   *mongo.Database
	collection *mongo.Collection
}

type authService interface {
	Login(user *models.User) ([]models.User, error)
	Register(user *models.User) (models.User, error)
}

func NewAuthService() (*AuthService, error) {
	client, err := config.GetClient()
	DataBase := client.Database("test")
	_collection := DataBase.Collection("users")
	// _collection.Drop(context.TODO()) //uncomment this tho drop collection
	if err == nil {
		return &AuthService{
			validator:  validator.New(),
			client:     client,
			DataBase:   DataBase,
			collection: _collection,
		}, nil
	} else {
		return nil, err
	}
}

func (as *AuthService) Login(user *models.User) (string, error, int) {

	// TODO: Implement user login logic
	filter := bson.D{{"email", user.Email}}
	var existingUser models.User
	err := as.collection.FindOne(context.TODO(), filter).Decode(&existingUser)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)) != nil {
		return "", errors.New("Invalid credentials"), http.StatusNotFound
	}

	// Generate JWT
	jwtToken, err := SignJwt(existingUser)
	if err != nil {
		return "", err, 500
	}

	return jwtToken, nil, 200
}

func (as *AuthService) Register(user *models.User) (models.OmitedUser, error, int) {
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

func SignJwt(existingUser models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  existingUser.ID,
		"email":    existingUser.Email,
		"is_admin": existingUser.Is_Admin,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
		// for the purpose of this task, we will set the expiration time to 1 minute
	})

	jwtToken, err := token.SignedString(config.JwtSecret)
	return jwtToken, err
}
