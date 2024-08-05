package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variable to store the JWT secret
var JwtSecret = []byte("your_jwt_secret")
var _client *mongo.Client
var _err error = nil

func MongoInit() (*mongo.Client, error) {
	// Read MONGO_CONNECTION_STRING from environment
	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING")

	clientOptions := options.Client().ApplyURI(MONGO_CONNECTION_STRING)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		_err = err
		log.Fatal(err)
		return nil, errors.New("Failed to connect to MongoDB")
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		_err = err
		log.Fatal(err)
		return nil, errors.New("Failed to connect to MongoDB")
	}

	fmt.Println("Connected to MongoDB!")
	_client = client
	return client, nil
}

func GetClient() (*mongo.Client, error) {
	return _client, _err
}
