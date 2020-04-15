package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomClient struct {
	client *mongo.Client
}

const (
	dbUrl  = "mongodb://localhost:27017"
	dbName = "test_task"
)

func connectToDb() (CustomClient, error) {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err :=
		mongo.Connect(context.TODO(), clientOptions)
	customClient := CustomClient{client: client}
	err = customClient.pingDataBase()
	return customClient, err
}

func (c *CustomClient) pingDataBase() error {
	err := c.client.Ping(context.TODO(), nil)
	return err
}

func (c *CustomClient) disconnectFromDb() error {
	err := c.client.Disconnect(context.TODO())
	return err
}

func (c *CustomClient) getCollection(name string) *mongo.Collection {
	return c.client.Database(dbName).Collection(name)
}

func insertOneToCollection(col mongo.Collection, value interface{}) (*mongo.InsertOneResult, error) {
	return col.InsertOne(context.TODO(), value)
}

func findOneById(col mongo.Collection, id primitive.ObjectID, elem interface{}) error {
	filter := bson.D{{"_id", id}}
	err := col.FindOne(context.TODO(), filter).Decode(elem)
	return err
}

func deleteFromDb(id interface{}, collection mongo.Collection) {
	id, _ = doPrettyId(fmt.Sprint(id))
	filter := bson.M{"_id": id}
	collection.DeleteOne(context.TODO(), filter)
}
