package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/url"
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
	user, err := h.client.getUserFromSession(r)

	if err != nil {
		fmt.Fprintf(w, "can't find user from the session : %v", err)
		return
	}

	h.user = user

	if strings.Contains(r.URL.Path, "createDoc") {
		h.createDocument(w, r)
	} else if strings.Contains(r.URL.Path, "editDoc") {
		h.editDocument(w, r)
	} else {
		h.deleteDocument(w, r)
	}
}

func (h *docHandler) createDocument(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == http.MethodPost {
		err = h.create(r.FormValue("Title"), r.FormValue("Content"))
		if err == nil {
			http.Redirect(w, r, "/", 302)
		}
	}

	//create inputs for form
	documentInput := getDocumentInput(nil)

	//execute template with data
	h.executeDocTemplate(w, documentInput, err, "Add a new document")
}

func (h *docHandler) create(title string, content string) error {
	docs := h.user.Documents

	doc := Document{
		Title:   title,
		Content: content,
	}
	insertedResult, err := h.insertDocument(doc)
	if err != nil {
		err = fmt.Errorf("can't insert a document: %v", err)
		return err
	}
	docs = append(docs, insertedResult.InsertedID)
	return h.updateUserDocs(docs)
}

func (h *docHandler) editDocument(w http.ResponseWriter, r *http.Request) {
	var err error

	doc, err := h.getDocument(r.URL)
	if err == nil && r.Method == http.MethodPost {
		err = h.edit(r.FormValue("Title"), r.FormValue("Content"), doc)
		if err == nil {
			http.Redirect(w, r, "/", 302)
		}
	}

	//create inputs for form
	documentInput := getDocumentInput(doc)

	//execute template with data
	h.executeDocTemplate(w, documentInput, err, "Edit the document")
}

func (h *docHandler) edit(title string, content string, document *Document) error {
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "title", Value: title},
			primitive.E{Key: "content", Value: content},
		}},
	}
	_, err := h.docCol.UpdateOne(context.TODO(), document, update)
	return err
}

func (h *docHandler) executeDocTemplate(w http.ResponseWriter, input *DocumentInput, err error, title string) {
	executeTemplate("views/createEdit.html", w, struct {
		Title         string
		DocumentInput *DocumentInput
		Error         error
		User          User
	}{
		Title:         title,
		DocumentInput: input,
		Error:         err,
		User:          *h.user,
	})
}

func (h *docHandler) deleteDocument(w http.ResponseWriter, r *http.Request) {
	docId := r.URL.Query().Get("docId")
	id, err := doPrettyId(fmt.Sprint(docId))
	if err != nil {
		fmt.Fprintf(w, "can't do normal id for search %v : %v", docId, err)
		return
	}
	err = h.deleteDoc(id)
	if err != nil {
		fmt.Fprintf(w, "can't delete a doc with id  %v : %v", id, err)
		return
	}
	update := bson.D{
		primitive.E{Key: "$pull", Value: bson.D{
			primitive.E{Key: "documents", Value: id},
		}}}

	err = h.updateUser(update)
	if err != nil {
		fmt.Fprintf(w, "can't delete a doc from user id  %v : %v", id, err)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func (h *docHandler) deleteDoc(id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := h.docCol.DeleteOne(context.TODO(), filter)
	return err
}

func (h *docHandler) updateUserDocs(docs []interface{}) error {
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "documents", Value: docs},
		}},
	}
	return h.updateUser(update)
}

func (h *docHandler) updateUser(update primitive.D) error {
	userCol := h.client.getCollection(userColName)
	_, err := userCol.UpdateOne(context.TODO(), h.user, update)
	return err
}

func (h *docHandler) getDocument(url *url.URL) (*Document, error) {
	id := fmt.Sprint(url.Query().Get("docId"))
	objectId, err := doPrettyId(id)

	if err != nil {
		return nil, fmt.Errorf("can't do normal id for search %v : %v", id, err)
	}

	var doc Document

	err = h.client.findOneById(docColName, objectId, &doc)

	if err != nil {
		return nil, fmt.Errorf("can't find a doc with id %v : %v", objectId, err)
	}

	return &doc, nil
}

func getDocumentInput(doc *Document) *DocumentInput {
	var inputs []Input

	title := ""
	content := ""

	if doc != nil {
		title = doc.Title
		content = doc.Content
	}

	input := Input{Name: "Title", Caption: "Enter title", Type: "input", Value: title, Required: true}
	inputs = append(inputs, input)
	input2 := Input{Name: "Content", Caption: "Enter content", Type: "textarea", Value: content, Required: true}
	inputs = append(inputs, input2)

	return &DocumentInput{inputs}
}

func (h *docHandler) insertDocument(value interface{}) (*mongo.InsertOneResult, error) {
	return h.docCol.InsertOne(context.TODO(), value)
}
