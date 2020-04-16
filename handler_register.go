package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		user, err := getUserByLogin(login, *collection)
		if err != nil {
			createErrorAndAppendToSlice(errors, err.Error())
		}

		if reflect.DeepEqual(user, User{}) {
			hash, error := HashPassword(r.FormValue("password"))
			if error == nil {
				user = User{
					Login:    login,
					Password: hash,
				}
				_, err := insertOneToCollection(*collection, user)
				if err != nil {
					fmt.Println("correct it")
				}
				http.Redirect(w, r, "/authorization", 302)
			} else {
				createErrorAndAppendToSlice(errors, "User is not registered. Try again!")
			}
		} else {
			createErrorAndAppendToSlice(errors, "User's already existed!")
		}

		if len(errors) != 0 {
			registerData.Errors = errors
		}
	}
	err := tmpl.Execute(w, registerData)
	if err != nil {
		fmt.Println("fwerhf")
	}
}

func createErrorAndAppendToSlice(errors []Error, name string) []Error {
	errors = append(errors, Error{
		Name: name,
	})
	return errors
}

func getUserByLogin(login string, collection mongo.Collection) (User, error) {
	var user User
	filter := bson.D{primitive.E{Key: "login", Value: login}}
	error := collection.FindOne(context.TODO(), filter).Decode(&user)
	return user, error
}
