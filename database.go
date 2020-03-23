package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func connectToDb() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("Cannot connect to db")
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Cannot listen db")
	}

	log.Println("Connected to MongoDB!")
	return client
}

func disconnectFromDb() {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal("Cannot disconnect")
	}
	log.Println("Disconnected from MongoDB!")
}

func getNeccessaryCollections(name string) *mongo.Collection {
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
