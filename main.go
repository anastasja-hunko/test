package main

import (
	"net/http"
)

type Error struct {
	Name string
}

type Document struct {
	User    User
	Name    string
	Content string
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/autorization", authorization)
	http.HandleFunc("/logout", logout)

	http.ListenAndServe("localhost:8181", nil)
}

//registration and storage of users
//conversation with outside world
//user can create his Documents
