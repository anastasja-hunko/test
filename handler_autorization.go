package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
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

//big piece, should be divided
func (h *authoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Println("ccf")
	}

	tmpl := template.Must(template.ParseFiles("views/userForm.html"))

	registerData := UserPostData{
		PageTitle: "Authorization",
	}

	if r.Method == http.MethodPost {
		var collection = h.client.getCollection("users")
		var errors []Error

		login := r.FormValue("login")
		user, _ := getUserByLogin(login, *collection)
		if reflect.DeepEqual(user, User{}) {
			errors = append(errors, Error{
				Name: "User is absent in database",
			})
		} else {
			password := r.FormValue("password")
			if CheckPasswordHash(password, user.Password) {
				session.Values["authorize"] = true
				session.Values["login"] = login
				err = sessions.Save(r, w)
				if err == nil {
					http.Redirect(w, r, "/", 302)
				}
			} else {
				errors = append(errors, Error{
					Name: "Incorrect password",
				})
			}
		}

		if len(errors) != 0 {
			registerData.Errors = errors
		}
	}
	tmpl.Execute(w, registerData)

}
