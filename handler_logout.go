package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")

	if err != nil {
		//obrabotat'
	}

	session.Values["authorize"] = false
	session.Values["login"] = ""
	err = sessions.Save(r, w)
	if err == nil {
		http.Redirect(w, r, "/", 302)
	}
}
