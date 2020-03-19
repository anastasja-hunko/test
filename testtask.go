package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type User struct {
	Login    string
	Password string
}

type UserPostData struct {
	PageTitle string
	Success   bool
	Errors    []Error
}

type Error struct {
	Name string
}

type Document struct {
	User    User
	Name    string
	Content string
}

func index(w http.ResponseWriter, r *http.Request) {
}

func register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("userForm.html"))

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

func getNeccessaryCollections(name string, client mongo.Client) *mongo.Collection {
	return client.Database("test_task").Collection(name)
}

func insertOneToCollection(col mongo.Collection, value interface{}, errors []Error) {
	insertResult, err := col.InsertOne(context.TODO(), value)

	if err != nil {
		errors = append(errors, Error{
			Name: "Cannot insert To Db",
		})
	}

	fmt.Println("Insertes one!", insertResult.InsertedID)
}

func connectToDb(errors []Error) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		errors = append(errors, Error{
			Name: "Cannot connect to db",
		})
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		errors = append(errors, Error{
			Name: "Cannot listen db",
		})
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func authorization(w http.ResponseWriter, r *http.Request) {
}

func logout(w http.ResponseWriter, r *http.Request) {
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
