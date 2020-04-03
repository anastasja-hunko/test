package main

import (
	"net/http"
)

type Error struct {
	Name string
}

var client = connectToDb()

func main() {
	defer disconnectFromDb()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/authorization", authorization)
	http.HandleFunc("/createDoc/", createDocument)
	http.HandleFunc("/editDoc/", editDocument)
	http.HandleFunc("/deleteDoc/", deleteDocument)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe("localhost:8181", nil)
}
