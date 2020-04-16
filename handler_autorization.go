package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"reflect"
)

var store = sessions.NewCookieStore([]byte("very-secret-key"))

type authoHandler struct {
	client *CustomClient
}

func newAuthoHandler(client *CustomClient) *authoHandler {
	return &authoHandler{client: client}
}

func (h *authoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resultErrors []error

	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		resultErrors = h.authoriseUser(resultErrors, login, r.FormValue("password"))
		if len(resultErrors) == 0 {
			err := workWithSession(w, r, true, login)
			if err == nil {
				http.Redirect(w, r, "/", 302)
			}
			resultErrors = append(resultErrors, err)
		}
	}

	//execute template with data
	executeTemplate("views/userForm.html", w, struct {
		PageTitle string
		Errors    []error
	}{
		PageTitle: "Authorization",
		Errors:    resultErrors,
	})
}

func (h *authoHandler) authoriseUser(resultErrors []error, login string, password string) []error {
	user, err := h.client.getUserByLogin(login)
	if err != nil || reflect.DeepEqual(user, User{}) {
		err = errors.New("user is absent in database:" + login)
		resultErrors = append(resultErrors, err)
		return resultErrors
	}

	if !CheckPasswordHash(password, user.Password) {
		err = errors.New("incorrect password")
		resultErrors = append(resultErrors, err)
		return resultErrors
	}
	return resultErrors
}

func workWithSession(w http.ResponseWriter, r *http.Request, authorize bool, login string) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		err = fmt.Errorf("can't get session with name %v", sessionName)
		return err
	}
	session.Values[sessionAuthorizeKey] = authorize
	session.Values[sessionLoginKey] = login
	return sessions.Save(r, w)
}
