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

// im not using this shiyt

func createDatabase(client *mongo.Client, dbName string) (*mongo.Database, error) {
	return client.Database(dbName), nil
}
func createCollection(db *mongo.Database, collectionName string) (*mongo.Collection, error) {
	return db.Collection(collectionName), nil
}

// crud start
func insertOne(collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(context.Background(), document)
}
func insertMany(collection *mongo.Collection, documents []interface{}) (*mongo.InsertManyResult, error) {
	return collection.InsertMany(context.Background(), documents)
}
func findOne(collection *mongo.Collection, filter interface{}) *mongo.SingleResult {
	return collection.FindOne(context.Background(), filter)
}
func find(collection *mongo.Collection, filter interface{}) (*mongo.Cursor, error) {
	return collection.Find(context.Background(), filter)
}
func updateOne(collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(context.Background(), filter, update)
}
func updateMany(collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return collection.UpdateMany(context.Background(), filter, update)
}
func deleteOne(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(context.Background(), filter)
}
func deleteMany(collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	return collection.DeleteMany(context.Background(), filter)
}

// crud end

func dropCollection(collection *mongo.Collection) error {
	return collection.Drop(context.Background())
}
func dropDatabase(db *mongo.Database) error {
	return db.Drop(context.Background())
}
func countDocuments(collection *mongo.Collection, filter interface{}) (int64, error) {
	return collection.CountDocuments(context.Background(), filter)
}
func listDatabases(client *mongo.Client) (mongo.ListDatabasesResult, error) {
	return client.ListDatabases(context.Background(), nil)
}
func listCollections(db *mongo.Database) ([]string, error) {
	return db.ListCollectionNames(context.Background(), nil)
}

// probably won't use this one
func aggregate(collection *mongo.Collection, pipeline interface{}) (*mongo.Cursor, error) {
	return collection.Aggregate(context.Background(), pipeline)
}

func createIndex(collection *mongo.Collection, model mongo.IndexModel) (string, error) {
	return collection.Indexes().CreateOne(context.Background(), model)
}

func dropIndex(collection *mongo.Collection, name string) error {
	_, err := collection.Indexes().DropOne(context.Background(), name)
	return err
}

func dropAllIndexes(collection *mongo.Collection) error {
	_, err := collection.Indexes().DropAll(context.Background())
	return err
}

func createUniqueIndex(collection *mongo.Collection, model mongo.IndexModel) (string, error) {
	model.Options = options.Index().SetUnique(true)
	return collection.Indexes().CreateOne(context.Background(), model)
}
