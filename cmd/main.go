package main

import (
	"github.com/anastasja-hunko/test/internal"
	"github.com/gorilla/sessions"
	"log"
)

func main() {
	config := internal.NewConfig()

	sessionStore := sessions.NewCookieStore([]byte("very-secret-key"))

	server := internal.New(config, sessionStore)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	//
	////file server for static content: js, css
	//fs := http.FileServer(http.Dir("static"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	//
	////endpoints
	//http.Handle("/", newIndexHandler(client))
	//
	//docHandler := newDocHandler(client)
	//http.Handle("/createDoc/", docHandler)
	//http.Handle("/editDoc/", docHandler)
	//http.Handle("/deleteDoc/", docHandler)
	//http.HandleFunc("/logout", logout)
	//
	////init and listen server
	//err = http.ListenAndServe(serverUrl, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
