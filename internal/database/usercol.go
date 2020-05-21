package database

import (
	"context"
	"github.com/anastasja-hunko/test/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCol struct {
	col *mongo.Collection
}

func (db *Database) NewUserCol() *UserCol {
	return &UserCol{col: db.db.Collection(db.config.UserColName)}
}

func (uc *UserCol) Create(u *model.User) error {
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	_, err := uc.col.InsertOne(context.TODO(), u)

	if err != nil {
		return err
	}

	return nil
}

func (uc *UserCol) FindByLogin(login string) (*model.User, error) {
	filter := bson.D{primitive.E{Key: "login", Value: login}}
	var user model.User
	err := uc.col.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (uc *UserCol) SetUserDocs(u *model.User, docs []interface{}) error {
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "documents", Value: docs},
		}},
	}
	return uc.updateUser(u, update)
}

func (uc *UserCol) RemoveIdFromUserDocs(u *model.User, id primitive.ObjectID) error {
	update := bson.D{
		primitive.E{Key: "$pull", Value: bson.D{
			primitive.E{Key: "documents", Value: id},
		}},
	}
	return uc.updateUser(u, update)
}

func (uc *UserCol) updateUser(u *model.User, update primitive.D) error {
	_, err := uc.col.UpdateOne(context.TODO(), u, update)
	return err
}
