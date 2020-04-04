package main

import (
	"net/http"
)

type Error struct {
	Name string
}

//!так нельзя делать
var client = connectToDb()

func main() {
	defer disconnectFromDb()
	//клиент и ошибку = коннект ту бд. + обработка ошибок

	//file server for static content: js, css
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//endpoints (add comments here. jpen api may
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/authorization", authorization)
	http.HandleFunc("/createDoc/", createDocument)
	http.HandleFunc("/editDoc/", editDocument)
	http.HandleFunc("/deleteDoc/", deleteDocument)
	http.HandleFunc("/logout", logout)

	//init and listen server
	http.ListenAndServe("localhost:8181", nil)
}
