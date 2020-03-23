package main

import (
	"html/template"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Println("ccf")
	}

	var tmpl *template.Template
	if session.Values["authorize"] == true {
		login := session.Values["login"]
		tmpl = template.Must(template.ParseFiles("views/authorizeMain.html"))
		tmpl.Execute(w, struct {
			Login interface{}
		}{
			Login: login,
		})
	} else {
		tmpl = template.Must(template.ParseFiles("views/NonMain.html"))
		tmpl.Execute(w, nil)
	}

	// область для документов
	// погода или курс валют
}
