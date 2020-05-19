package main

import (
	"fmt"
	"log"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	err := workWithSession(w, r, false, "")
	if err != nil {
		log.Fatal(err)
		fmt.Fprint(w, err)
	}
	http.Redirect(w, r, "/", 302)
}
