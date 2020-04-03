package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
		log.Println("ccf")
	}

	var tmpl *template.Template
	if session.Values["authorize"] == true {
		tmpl = template.Must(template.ParseFiles("views/indexWhenAuthorized.html"))
		login := session.Values["login"]

		var collection = getNeccessaryCollections("users")
		user := checkLoginOnExisting(fmt.Sprintf("%v", login), *collection)

		var docCol = getNeccessaryCollections("docs")
		var documents []Document
		documents = getDocumentsByUser(user, *docCol)

		url := "http://www.nbrb.by/api/exrates/rates?periodicity=0"

		client := http.Client{
			Timeout: time.Second * 2,
		}

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		//ошибку проверить

		res, _ := client.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		var course []Course
		_ = json.Unmarshal(body, &course)

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

	filter := bson.D{{"userid", user.Id}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {

		var elem Document
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		docs = append(docs, elem)
	}
	return docs
}
