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

func connectToDb() (*CustomClient, error) {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, fmt.Errorf("Can't connect to database %v: %v ", dbUrl, err)
	}

	customClient := CustomClient{client: client}
	err = customClient.pingDataBase()

	if err != nil {
		return nil, fmt.Errorf("Can't ping the database %v: %v ", dbUrl, err)
	}

	return &customClient, err
}

func (c *CustomClient) pingDataBase() error {
	return c.client.Ping(context.TODO(), nil)
}

func (c *CustomClient) disconnectFromDb() error {
	return c.client.Disconnect(context.TODO())
}

func (c *CustomClient) getCollection(name string) *mongo.Collection {
	return c.client.Database(dbName).Collection(name)
}

func (c *CustomClient) getUserByLogin(login string) (*User, error) {
	filter := bson.D{primitive.E{Key: "login", Value: login}}
	collection := c.getCollection(userColName)

	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *CustomClient) insertOneToCollection(colName string, value interface{}) (*mongo.InsertOneResult, error) {
	col := c.getCollection(colName)
	return col.InsertOne(context.TODO(), value)
}

func (c *CustomClient) findOneById(colName string, id primitive.ObjectID, elem interface{}) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	col := c.getCollection(colName)
	return col.FindOne(context.TODO(), filter).Decode(elem)
}
