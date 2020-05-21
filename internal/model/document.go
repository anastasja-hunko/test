package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	Title   string `bson:"title,"`
	Content string `bson:"content,"`
	Id      primitive.ObjectID
}

type Input struct {
	Name     string
	Caption  string
	Value    string
	Type     string
	Required bool
}

type DocumentInput struct {
	Inputs *[]Input
	Create bool
}

type DocView struct {
	DocumentInput *DocumentInput
	User          *User
}
