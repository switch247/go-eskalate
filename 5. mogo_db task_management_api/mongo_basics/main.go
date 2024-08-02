package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// Read MONGO_CONNECTION_STRING from environment
	MONGO_CONNECTION_STRING := os.Getenv("MONGO_CONNECTION_STRING")
	// Set client options
	clientOptions := options.Client().ApplyURI(MONGO_CONNECTION_STRING)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	DataBase := client.Database("test")
	collection := DataBase.Collection("trainers")

	// insert one
	InsertOne(collection)
	// insert many at once
	InsertMany(collection)

	// update one
	UpdateOne(collection)
	GetOne(collection)

	DeleteOne(collection)

	GetOne(collection)

	val, err := GetMany(collection)
	if err == nil {

		for _, t := range val {
			fmt.Println(t.Name)
		}
	}
	// close connection
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}

func GetOne(collection *mongo.Collection) {
	// create a value into which the result can be decoded
	// where name== Ash
	var result Trainer
	filter := bson.D{{"name", "Ash"}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
}

func GetMany(collection *mongo.Collection) ([]*Trainer, error) {
	// Pass these options to the Find method
	findOptions := options.Find()
	// findOptions.SetLimit(2)
	filter := bson.D{{}}

	// Here's an array in which you can store the decoded documents
	var results []*Trainer

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return []*Trainer{}, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil
}

func UpdateOne(collection *mongo.Collection) {
	// i hate this
	filter := bson.D{{"name", "Ash"}}
	// where name== Ash
	// age+=1
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func UpdateMany(collection *mongo.Collection) {
	// i hate this
	filter := bson.D{{"name", "Ash"}}
	// where name== Ash
	// age+=1
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func InsertOne(collection *mongo.Collection) {

	ash := Trainer{"Ash", 10, "Pallet Town"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func InsertMany(collection *mongo.Collection) {
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func DeleteOne(collection *mongo.Collection) {
	filter := bson.D{{"name", "Ash"}}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func DeleteMany(collection *mongo.Collection) {
	filter := bson.D{{}}
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func NukeACollection(collection *mongo.Collection) {
	err := collection.Drop(context.TODO())
	if err != nil {
		fmt.Println("Collection Dropped")
	} else {
		fmt.Println("Could Not Drop Collection:")
		log.Fatal(err)
	}
}

// D: A BSON document. This type should be used in situations where order matters, such as MongoDB commands.
// M: An unordered map. It is the same as D, except it does not preserve order.
// A: A BSON array.
// E: A single element inside a D.
// Here is an example of a filter document built using D types which may be used to find documents where the name field matches either Alice or Bob:

// bson.D{{
//     "name",
//     bson.D{{
//         "$in",
//         bson.A{"Alice", "Bob"}
//     }}
// }}
