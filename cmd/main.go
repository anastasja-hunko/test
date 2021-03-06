package main

import (
	"github.com/anastasja-hunko/test/internal"
	"github.com/gorilla/sessions"
	"log"
)

func main() {
	//init config, sessionStore and server. Then start server
	config := internal.NewConfig()

	sessionStore := sessions.NewCookieStore([]byte("very-secret-key"))

	server := internal.New(config, sessionStore)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
