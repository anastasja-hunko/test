package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToDb(errors []Error) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		errors = append(errors, Error{
			Name: "Cannot connect to db",
		})
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		errors = append(errors, Error{
			Name: "Cannot listen db",
		})
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func getNeccessaryCollections(name string, client mongo.Client) *mongo.Collection {
	return client.Database("test_task").Collection(name)
}

func insertOneToCollection(col mongo.Collection, value interface{}, errors []Error) {
	insertResult, err := col.InsertOne(context.TODO(), value)

	if err != nil {
		errors = append(errors, Error{
			Name: "Cannot insert To Db",
		})
	}

	fmt.Println("Insertes one!", insertResult.InsertedID)
}
