package main

import (
	"log"
	"net/http"
)

type Error struct {
	Name string
}

const (
	serverUrl = "localhost:8181"
)

func main() {
	client, err := connectToDb()
	if err != nil {
		log.Println(err, "can't connect to database")
	}

	defer func() {
		err := client.disconnectFromDb()
		log.Println(err.Error())
	}()

	//file server for static content: js, css
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//endpoints (add comments here. jpen api may
	http.Handle("/", newIndexHandler(&client))
	http.Handle("/register", newRegisterHandler(&client))
	http.Handle("/authorization", newAuthoHandler(&client))

	docHandler := newDocHandler(&client)
	http.Handle("/createDoc/", docHandler)
	http.Handle("/editDoc/", docHandler)
	http.Handle("/deleteDoc/", docHandler)
	http.HandleFunc("/logout", logout)

	//init and listen server
	http.ListenAndServe(serverUrl, nil)
}
