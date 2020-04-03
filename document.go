package main

type Document struct {
	Title   string
	Content string
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
