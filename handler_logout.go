package main

import (
	"fmt"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	err := workWithSession(w, r, false, "")
	if err != nil {
		fmt.Fprint(w, err)
	}
	http.Redirect(w, r, "/", 302)
}
