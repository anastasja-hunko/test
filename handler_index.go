package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Course struct {
	Abbreviation string  `json:"Cur_Abbreviation"`
	Rate         float64 `json:"Cur_OfficialRate"`
}

type indexHandler struct {
	client *CustomClient
}

func newIndexHandler(client *CustomClient) *indexHandler {
	return &indexHandler{client: client}
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Println(err)
	}

	var tmpl *template.Template
	if session.Values["authorize"] == true {
		tmpl = template.Must(template.ParseFiles("views/indexWhenAuthorized.html"))
		login := session.Values["login"]

		var collection = h.client.getCollection("users")
		user, _ := getUserByLogin(fmt.Sprintf("%v", login), *collection)

		var docCol = h.client.getCollection("docs")
		var documents []Document
		documents, _ = getDocumentsByUser(user, *docCol)

		url := "http://www.nbrb.by/api/exrates/rates?periodicity=0"

		client := http.Client{
			Timeout: time.Second * 2,
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		var course []Course

		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		res, err := client.Do(req)

		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			err = json.Unmarshal(body, &course)
		} else {
			http.Error(w, err.Error(), 500)
		}

		tmpl.Execute(w, struct {
			User      User
			Course    []Course
			Documents []Document
		}{
			User:      user,
			Course:    course,
			Documents: documents,
		})
	} else {
		tmpl := template.Must(template.ParseFiles("views/indexWhenNonAuthorized.html"))
		tmpl.Execute(w, nil)
	}
}

func getDocumentsByUser(user User, collection mongo.Collection) ([]Document, error) {
	var docs []Document

	for d := range user.Documents {
		var elem Document
		id, err := doPrettyId(fmt.Sprint(user.Documents[d]))
		if err != nil {
			return docs, err
		}
		err = findOneById(collection, id, &elem)
		if err == nil {
			elem.Id = fmt.Sprint(user.Documents[d])
			if err != nil {
				return docs, err
			}
			docs = append(docs, elem)
		}
	}
	return docs, nil
}

func doPrettyId(stringId string) (primitive.ObjectID, error) {
	stringId = stringId[10 : len(stringId)-2]
	return primitive.ObjectIDFromHex(stringId)
}
