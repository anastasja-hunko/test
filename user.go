package main

import (
	"fmt"
	"net/http"
)

type User struct {
	Login     string
	Password  string
	Documents []interface{}
}

func (c *CustomClient) getUserFromSession(r *http.Request) (User, error) {
	session, err := store.Get(r, sessionName)
	var user User
	if err != nil {
		return user, err
	}
	login := session.Values[sessionLoginKey]
	user, err = c.getUserByLogin(fmt.Sprint(login))
	return user, err
}
