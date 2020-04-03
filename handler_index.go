package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

func index(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Println(err)
	}

	var tmpl *template.Template
	if session.Values["authorize"] == true {
		tmpl = template.Must(template.ParseFiles("views/indexWhenAuthorized.html"))
		login := session.Values["login"]

		var collection = getNeccessaryCollections("users")
		user := getUserByLogin(fmt.Sprintf("%v", login), *collection)

		var docCol = getNeccessaryCollections("docs")
		var documents []Document
		documents = getDocumentsByUser(user, *docCol)

		url := "http://www.nbrb.by/api/exrates/rates?periodicity=0"

		client := http.Client{
			Timeout: time.Second * 2,
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		var course []Course

		if err != nil {
			log.Println(err)
		} else {
			res, err := client.Do(req)
			if err != nil {
				log.Println(err)
			} else {
				body, err := ioutil.ReadAll(res.Body)
				if err == nil {
					err = json.Unmarshal(body, &course)
				}
			}
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

func getDocumentsByUser(user User, collection mongo.Collection) []Document {
	var docs []Document

	for d := range user.Documents {
		var stringId = fmt.Sprint(user.Documents[d])
		stringId = stringId[10 : len(stringId)-2]
		ff, _ := primitive.ObjectIDFromHex(stringId)
		filter := bson.D{{"_id", ff}}
		var elem Document
		err := collection.FindOne(context.TODO(), filter).Decode(&elem)
		elem.Id = stringId
		if err == nil {
			docs = append(docs, elem)
		}
	}

	return docs
}
