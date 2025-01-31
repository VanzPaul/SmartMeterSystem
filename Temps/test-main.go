package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string
	Age  int
	Sex  string
}

func login() Person {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://vanspaul09:ab7vSvvo14nx7gN3@cluster0.euhiz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Insert data to the database
	/*
		coll := client.Database("test_db").Collection("persons")
		newPerson := Person{
			Name: "John Doe",
			Age:  36,
			Sex:  "male",
		}
		result, err := coll.InsertOne(context.TODO(), newPerson)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Successfully inserted:\nName: %s\nAge: %d\nSex: %s\n", newPerson.Name, newPerson.Age, newPerson.Sex)
		fmt.Println("Result:")
		fmt.Println(result)
	*/

	// Get the data from the database
	coll := client.Database("test_db").Collection("persons")
	filter := bson.D{{Key: "name", Value: "John Doe"}}

	var result Person
	err = coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(err)
			panic(err)
		}
		panic(err)
	}
	//fmt.Println("Result:")
	//fmt.Printf("ID: %s\nName: %s\nAge: %d\nSex: %s\n", result.ID.Hex(), result.Name, result.Age, result.Sex)

	return result
}

func main() {

	var result Person

	result = login()

	fmt.Println("Result:")
	fmt.Printf("ID: %s\nName: %s\nAge: %d\nSex: %s\n", result.ID.Hex(), result.Name, result.Age, result.Sex)
}
