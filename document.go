package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string
	Content string
	UserId  interface{}
}

type Input struct {
	Name     string
	Caption  string
	Value    string
	Type     string
	Required bool
}

type DocumentInput struct {
	Inputs []Input
	Title  string
	User   User
}
