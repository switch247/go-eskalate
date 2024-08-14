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

// im not using this shiyt

// func createDatabase(client *mongo.Client, dbName string) (*mongo.Database, error) {
// 	return client.Database(dbName), nil
// }
// func createCollection(db *mongo.Database, collectionName string) (*mongo.Collection, error) {
// 	return db.Collection(collectionName), nil
// }

// type Collections interface {
// 	insertOne(document interface{}) (*mongo.InsertOneResult, error)
// 	insertMany(documents []interface{}) (*mongo.InsertManyResult, error)
// 	findOne(filter interface{}) *mongo.SingleResult
// 	find(filter interface{}) (*mongo.Cursor, error)
// 	updateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error)
// 	updateMany(filter interface{}, update interface{}) (*mongo.UpdateResult, error)
// 	deleteOne(filter interface{}) (*mongo.DeleteResult, error)
// 	deleteMany(filter interface{}) (*mongo.DeleteResult, error)
// 	dropCollection() error
// }

// type collection struct {
// }

// // crud start
// func (collection *collection) insertOne(document interface{}) (*mongo.InsertOneResult, error) {
// 	return collection.InsertOne(context.Background(), document)
// }
// func (collection *collection) insertMany(documents []interface{}) (*mongo.InsertManyResult, error) {
// 	return collection.InsertMany(context.Background(), documents)
// }
// func (collection *collection) findOne(filter interface{}) *mongo.SingleResult {
// 	return collection.FindOne(context.Background(), filter)
// }
// func (collection *collection) find(filter interface{}) (*mongo.Cursor, error) {
// 	return collection.Find(context.Background(), filter)
// }
// func (collection *collection) updateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
// 	return collection.UpdateOne(context.Background(), filter, update)
// }
// func (collection *collection) updateMany(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
// 	return collection.UpdateMany(context.Background(), filter, update)
// }
// func (collectio *collection) deleteOne(filter interface{}) (*mongo.DeleteResult, error) {
// 	return collection.DeleteOne(context.Background(), filter)
// }
// func (collection *Collections) deleteMany(filter interface{}) (*mongo.DeleteResult, error) {
// 	return collection.DeleteMany(context.Background(), filter)
// }

// // crud end

// func (collection *mongo.Collection) dropCollection() error {
// 	return collection.Drop(context.Background())
// }
// func  (collection *mongo.Collection) countDocuments( filter interface{}) (int64, error) {
// 	return collection.CountDocuments(context.Background(), filter)
// }

// // func dropDatabase(db *mongo.Database) error {
// // 	return db.Drop(context.Background())
// // }

// // func listDatabases(client *mongo.Client) (mongo.ListDatabasesResult, error) {
// // 	return client.ListDatabases(context.Background(), nil)
// // }
// // func listCollections(db *mongo.Database) ([]string, error) {
// // 	return db.ListCollectionNames(context.Background(), nil)
// // }

// // // probably won't use this one
// // func aggregate(collection *mongo.Collection, pipeline interface{}) (*mongo.Cursor, error) {
// // 	return collection.Aggregate(context.Background(), pipeline)
// // }

// // func createIndex(collection *mongo.Collection, model mongo.IndexModel) (string, error) {
// // 	return collection.Indexes().CreateOne(context.Background(), model)
// // }

// // func dropIndex(collection *mongo.Collection, name string) error {
// // 	_, err := collection.Indexes().DropOne(context.Background(), name)
// // 	return err
// // }

// // func dropAllIndexes(collection *mongo.Collection) error {
// // 	_, err := collection.Indexes().DropAll(context.Background())
// // 	return err
// // }

// // func createUniqueIndex(collection *mongo.Collection, model mongo.IndexModel) (string, error) {
// // 	model.Options = options.Index().SetUnique(true)
// // 	return collection.Indexes().CreateOne(context.Background(), model)
// // }
