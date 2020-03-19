package main

import (
	"html/template"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/userForm.html"))

	registerData := UserPostData{
		PageTitle: "Registration",
		Success:   false,
	}

	if r.Method == http.MethodPost {
		var errors []Error
		c := connectToDb(errors)

		//проверка на логин
		// проверка на пароль

		//получение данных

		hash, _ := HashPassword(r.FormValue("password"))

		user := User{
			Login:    r.FormValue("login"),
			Password: hash,
		}
		//сохранение в БД и сессию
		col := getNeccessaryCollections("users", *c)
		insertOneToCollection(*col, user, errors)

		if len(errors) == 0 {
			registerData = UserPostData{
				Success: true,
			}
		}
	}
	tmpl.Execute(w, registerData)

}
