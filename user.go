package main

import (
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Login     string
	Password  string
	Documents []interface{}
}

func (c *CustomClient) getUserFromSession(r *http.Request) User {
	session, err := store.Get(r, "session-name")
	var user User
	if err != nil {
		log.Println(err)
	} else {
		login := session.Values["login"]
		var collection = c.getCollection("users")
		user = getUserByLogin(fmt.Sprintf("%v", login), *collection)
	}
	return user
}
