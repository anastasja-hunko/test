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
	client  *mongo.Client
	context context.Context
}

func connectToDb() (CustomClient, error) {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		err = fmt.Errorf("Can't connect to database %v: %v ", dbUrl, err)
		return CustomClient{}, err
	}
	customClient := CustomClient{client: client, context: context.TODO()}
	err = customClient.pingDataBase()
	if err != nil {
		err = fmt.Errorf("Can't ping the database %v: %v ", dbUrl, err)
	}
	return customClient, err
}

func (c *CustomClient) pingDataBase() error {
	err := c.client.Ping(c.context, nil)
	return err
}

func (c *CustomClient) disconnectFromDb() error {
	err := c.client.Disconnect(c.context)
	return err
}

func (c *CustomClient) getCollection(name string) *mongo.Collection {
	return c.client.Database(dbName).Collection(name)
}

func (c *CustomClient) getUserByLogin(login string) (User, error) {
	var collection = c.getCollection(userColName)
	return getUserByLogin(login, *collection)
}

func insertOneToCollection(col mongo.Collection, value interface{}) (*mongo.InsertOneResult, error) {
	return col.InsertOne(context.TODO(), value)
}

func findOneById(col mongo.Collection, id primitive.ObjectID, elem interface{}) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	err := col.FindOne(context.TODO(), filter).Decode(elem)
	return err
}

func deleteFromDb(id interface{}, collection mongo.Collection) error {
	id, _ = doPrettyId(fmt.Sprint(id))
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}
