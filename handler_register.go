package main

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
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
	//execute template with data
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

	//check error
	if user != nil {
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
	user = &User{
		Login:    login,
		Password: hash,
	}
	_, err = h.insertUser(*user)

	if err != nil {
		err = fmt.Errorf("cannot insert user : %v", err)
		resultErrors = append(resultErrors, err)
		return resultErrors
	}
	return resultErrors
}

func (h *registerHandler) insertUser(user User) (*mongo.InsertOneResult, error) {
	return h.client.insertOneToCollection(userColName, user)
}
