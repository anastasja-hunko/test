package main

type Document struct {
	Title   string `bson:"title,"`
	Content string `bson:"content,"`
	Id      string
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
