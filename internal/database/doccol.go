package database

import (
	"context"
	"fmt"
	"github.com/anastasja-hunko/test/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocCol struct {
	col *mongo.Collection
}

func (db *Database) NewDocCol() *DocCol {
	return &DocCol{col: db.db.Collection(db.config.DocColName)}
}

func (dc *DocCol) FindById(id primitive.ObjectID) (*model.Document, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	var elem model.Document
	if err := dc.col.FindOne(context.TODO(), filter).Decode(&elem); err != nil {
		return nil, err
	}
	return &elem, nil
}

func (dc *DocCol) FindDocumentsByUser(user *model.User) (*[]model.Document, error) {
	var docs []model.Document

	for d := range user.Documents {
		id, err := doPrettyId(fmt.Sprint(user.Documents[d]))

		if err != nil {
			return &docs, fmt.Errorf("Can't do id for search in database %v: %v ", id, err)
		}

		elem, err := dc.FindById(id)
		if err != nil {
			return &docs, fmt.Errorf("Can't find document with id %v: %v ", id, err)
		}

		docs = append(docs, *elem)
	}
	return &docs, nil
}

func doPrettyId(stringId string) (primitive.ObjectID, error) {
	stringId = stringId[10 : len(stringId)-2]
	return primitive.ObjectIDFromHex(stringId)
}

func (dc *DocCol) Create(d *model.Document) (*mongo.InsertOneResult, error) {
	ir, err := dc.col.InsertOne(context.TODO(), d)

	if err != nil {
		return nil, err
	}

	return ir, nil
}

func (dc *DocCol) Edit(document *model.Document) error {
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "title", Value: document.Title},
			primitive.E{Key: "content", Value: document.Content},
		}},
	}
	_, err := dc.col.UpdateOne(context.TODO(), document, update)
	return err
}

func (dc *DocCol) Delete(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := dc.col.DeleteOne(context.TODO(), filter)
	return err
}
