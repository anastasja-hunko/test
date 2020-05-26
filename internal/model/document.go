package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//document struct
type Document struct {
	Title   string `bson:"title,"`
	Content string `bson:"content,"`
	Id      primitive.ObjectID
}

//input struct: use for html
type Input struct {
	Name     string
	Caption  string
	Value    string
	Type     string
	Required bool
}

//document input: use for documents forms
type DocumentInput struct {
	Inputs *[]Input
	Create bool
}

type DocView struct {
	DocumentInput *DocumentInput
	User          *User
}
