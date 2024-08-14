package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	mongo "main/mongo"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

var _client mongo.Client
var _err error = nil

func MongoInit() (mongo.Client, error) {
	Envinit()

	// clientOptions := options.Client().ApplyURI(MONGO_CONNECTION_STRING)

	client, err := mongo.NewClient(MONGO_CONNECTION_STRING)
	//  mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		_err = err
		log.Fatal(err)
		return nil, errors.New("Failed to connect to MongoDB")
	}

	err = client.Ping(context.TODO())

	if err != nil {
		_err = err
		log.Fatal(err)
		return nil, errors.New("Failed to connect to MongoDB")
	}

	fmt.Println("Connected to MongoDB!")
	_client = client
	return _client, nil
}

func GetClient() (mongo.Client, error) {
	return _client, _err
}
