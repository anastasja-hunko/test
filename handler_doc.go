package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"log"
	"net/http"
)

func createDocument(w http.ResponseWriter, r *http.Request) {
	userCol := getNeccessaryCollections("users")
	userLogin := r.URL.Query().Get("login")
	user := checkLoginOnExisting(userLogin, *userCol)

	docCol := getNeccessaryCollections("docs")
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("views/createEdit.html"))
		var inputs []Input

		input := Input{Name: "Title", Caption: "Enter title", Type: "input", Required: true}
		inputs = append(inputs, input)
		input2 := Input{Name: "Content", Caption: "Enter content", Type: "textarea", Required: true}
		inputs = append(inputs, input2)

		documentInput := DocumentInput{
			Inputs: inputs,
			Title:  "Add a new document!",
			User:   user,
		}
		tmpl.Execute(w, documentInput)
	} else {
		doc := Document{
			Title:   r.FormValue("Title"),
			Content: r.FormValue("Content"),
			UserId:  user.Id,
		}
		insertOneToCollection(*docCol, doc, []Error{})

		http.Redirect(w, r, "/", 302)
	}
}

func editDocument(w http.ResponseWriter, r *http.Request) {
	docId := r.URL.Query().Get("docId")
	docCol := getNeccessaryCollections("docs")

	doc := getDocById(docId, *docCol)

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("views/createEdit.html"))
		var inputs []Input

		input := Input{Name: "Title", Caption: "Enter title", Type: "input", Value: doc.Title, Required: true}
		inputs = append(inputs, input)
		input2 := Input{Name: "Content", Caption: "Enter content", Type: "textarea", Value: doc.Content, Required: true}
		inputs = append(inputs, input2)

		documentInput := DocumentInput{
			Inputs: inputs,
			Title:  "Edit the document!",
			User:   User{},
		}
		tmpl.Execute(w, documentInput)
	} else {
		update := bson.D{
			{"$set", bson.D{
				{"title", r.FormValue("Title")},
				{"content", r.FormValue("Content")},
			}},
		}

		_, err := docCol.UpdateOne(context.TODO(), doc, update)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", 302)
	}
}

func deleteDocument(w http.ResponseWriter, r *http.Request) {
	docId := r.URL.Query().Get("docId")
	docCol := getNeccessaryCollections("docs")
	deleteFromDb(docId, *docCol)
	http.Redirect(w, r, "/", 302)
}

func getDocById(id interface{}, collection mongo.Collection) Document {
	var doc Document
	stringId := fmt.Sprint(id)
	stringId = stringId[10 : len(stringId)-2]
	id, _ = primitive.ObjectIDFromHex(stringId)
	filter := bson.M{"_id": id}
	collection.FindOne(context.TODO(), filter).Decode(&doc)
	return doc
}

func deleteFromDb(id interface{}, collection mongo.Collection) {
	stringId := fmt.Sprint(id)
	stringId = stringId[10 : len(stringId)-2]
	id, _ = primitive.ObjectIDFromHex(stringId)
	filter := bson.M{"_id": id}
	collection.DeleteOne(context.TODO(), filter)
}
