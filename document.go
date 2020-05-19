package main

type Document struct {
	Title   string "json:title"
	Content string "json:content"
	Id      string "json:id"
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
}
