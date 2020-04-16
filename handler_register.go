package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"reflect"
)

type registerHandler struct {
	client *CustomClient
}

func newRegisterHandler(client *CustomClient) *registerHandler {
	return &registerHandler{client: client}
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errors []error

	if r.Method == http.MethodPost {
		errors = h.registerUser(r, errors)

		if len(errors) == 0 {
			http.Redirect(w, r, "/authorization", 302)
		}
	}

	executeTemplate("views/userForm.html", w, struct {
		PageTitle string
		Errors    []error
	}{
		PageTitle: "Registration",
		Errors:    errors,
	})
}

func (h *registerHandler) registerUser(r *http.Request, resultErrors []error) []error {
	login := r.FormValue("login")
	user, err := h.client.getUserByLogin(login)

	if err == nil && !reflect.DeepEqual(user, User{}) {
		err = errors.New("user's already existed with login:" + login)
		resultErrors = append(resultErrors, err)
		return resultErrors
	}
	hash, err := HashPassword(r.FormValue("password"))
	if err != nil {
		err = fmt.Errorf("cannot hash pasword %v : %v", r.FormValue("password"), err)
		resultErrors = append(resultErrors, err)
		return resultErrors
	}
	user = User{
		Login:    login,
		Password: hash,
	}
	_, err = h.insertUser(user)

	if err != nil {
		err = fmt.Errorf("cannot insert user : %v", err)
		resultErrors = append(resultErrors, err)
		return resultErrors
	}
	return resultErrors
}

func getUserByLogin(login string, collection mongo.Collection) (User, error) {
	var user User
	filter := bson.D{primitive.E{Key: "login", Value: login}}
	error := collection.FindOne(context.TODO(), filter).Decode(&user)
	return user, error
}

func (h *registerHandler) insertUser(user User) (*mongo.InsertOneResult, error) {
	var collection = h.client.getCollection(userColName)
	return insertOneToCollection(*collection, user)
}
