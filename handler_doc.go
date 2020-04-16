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
	"strings"
)

type docHandler struct {
	client *CustomClient
	docCol *mongo.Collection
	user   *User
}

func newDocHandler(client *CustomClient) *docHandler {
	return &docHandler{client: client, docCol: client.getCollection(docColName)}
}

func (h *docHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user := h.client.getUserFromSession(r)
	h.user = &user

	if strings.Contains(r.URL.Path, "createDoc") {
		h.createDocument(w, r)
	} else if strings.Contains(r.URL.Path, "editDoc") {
		h.editDocument(w, r)
	} else {
		h.deleteDocument(w, r)
	}
}

func (h *docHandler) createDocument(w http.ResponseWriter, r *http.Request) {
	userCol := h.client.getCollection(userColName)
	userLogin := r.URL.Query().Get("login")
	user, _ := getUserByLogin(userLogin, *userCol)

	if r.Method == http.MethodGet {
		h.showDocForm(w, Document{}, "Add a new document!")
	} else {
		doc := Document{
			Title:   r.FormValue("Title"),
			Content: r.FormValue("Content"),
		}
		insertedResult, err := insertOneToCollection(*h.docCol, doc)
		if err != nil {
			fmt.Println("correct it")
		}
		docs := user.Documents
		docs = append(docs, insertedResult.InsertedID)

		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "documents", Value: docs},
			}},
		}
		_, err = userCol.UpdateOne(context.TODO(), user, update)

		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", 302)
	}
}

func (h *docHandler) editDocument(w http.ResponseWriter, r *http.Request) {
	docId := r.URL.Query().Get("docId")
	objectId, err := doPrettyId(fmt.Sprint(docId))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	var doc Document

	err = findOneById(*h.docCol, objectId, &doc)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if r.Method == http.MethodGet {
		h.showDocForm(w, doc, "Edit the document")
	} else {
		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "title", Value: r.FormValue("Title")},
				primitive.E{Key: "content", Value: r.FormValue("Content")},
			}},
		}
		filter := bson.M{"_id": objectId}
		_, err := h.docCol.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", 302)
	}
}

func (h *docHandler) showDocForm(w http.ResponseWriter, doc Document, title string) {
	tmpl := template.Must(template.ParseFiles("views/createEdit.html"))
	var inputs []Input

	input := Input{Name: "Title", Caption: "Enter title", Type: "input", Value: doc.Title, Required: true}
	inputs = append(inputs, input)
	input2 := Input{Name: "Content", Caption: "Enter content", Type: "textarea", Value: doc.Content, Required: true}
	inputs = append(inputs, input2)

	documentInput := DocumentInput{inputs, title, *h.user}

	err := tmpl.Execute(w, documentInput)
	if err != nil {
		fmt.Println("correct it")
	}
}

func (h *docHandler) deleteDocument(w http.ResponseWriter, r *http.Request) {
	docId := r.URL.Query().Get("docId")
	err := deleteFromDb(docId, *h.docCol)
	if err != nil {
		fmt.Println("correct it")
	}
	http.Redirect(w, r, "/", 302)
}
