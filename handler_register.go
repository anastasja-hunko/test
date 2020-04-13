package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"reflect"
)

type registerHandler struct {
	client    *CustomClient
	pageTitle string
}

func newRegisterHandler(client *CustomClient) *registerHandler {
	return &registerHandler{client: client, pageTitle: "Registration"}
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/userForm.html"))

	registerData := UserPostData{
		PageTitle: h.pageTitle,
	}

	if r.Method == http.MethodPost {
		var collection = h.client.getCollection("users")
		var errors []Error

		login := r.FormValue("login")
		user := getUserByLogin(login, *collection)

		if reflect.DeepEqual(user, User{}) {
			hash, error := HashPassword(r.FormValue("password"))
			if error == nil {
				user = User{
					Login:    login,
					Password: hash,
				}
				insertOneToCollection(*collection, user)
				http.Redirect(w, r, "/authorization", 302)
			} else {
				http.Error(w, error.Error(), 500)
			}
		} else {
			errors = append(errors, Error{
				Name: "User's already exist!",
			})
		}

		if len(errors) != 0 {
			registerData.Errors = errors
		}
	}
	tmpl.Execute(w, registerData)
}

func getUserByLogin(login string, collection mongo.Collection) User {
	var user User

	filter := bson.D{{"login", login}}
	collection.FindOne(context.TODO(), filter).Decode(&user)
	return user
}
