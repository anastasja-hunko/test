package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"reflect"
)

func register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/userForm.html"))

	registerData := UserPostData{
		PageTitle: "Registration",
	}

	if r.Method == http.MethodPost {
		var collection = getNeccessaryCollections("users")
		var errors []Error

		login := r.FormValue("login")
		user := checkLoginOnExisting(login, *collection)

		if reflect.DeepEqual(user, User{}) {
			hash, _ := HashPassword(r.FormValue("password"))
			user = User{
				Login:    login,
				Password: hash,
			}
			insertOneToCollection(*collection, user, errors)
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

func checkLoginOnExisting(login string, collection mongo.Collection) User {
	var user User

	filter := bson.D{{"login", login}}
	collection.FindOne(context.TODO(), filter).Decode(&user)
	return user
}
