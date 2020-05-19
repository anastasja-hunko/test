package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Course struct {
	Abbreviation string  `json:"Cur_Abbreviation"`
	Rate         float64 `json:"Cur_OfficialRate"`
}

type indexHandler struct {
	client *CustomClient
}

func newIndexHandler(client *CustomClient) *indexHandler {
	return &indexHandler{client: client}
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)

	if err != nil || session.Values[sessionAuthorizeKey] != true {
		//if session is unavailable or doesn't consist authorizeKey then
		//show register and authorize link
		log.Println(err)
		executeTemplate("views/indexWhenNonAuthorized.html", w, nil)
		return
	}

	//get courser from nbrb
	course, err := getCourses()

	login := session.Values[sessionLoginKey]

	user, err2 := h.client.getUserByLogin(fmt.Sprint(login))
	documents, docErr := h.getDocumentsByUser(user)

	if docErr != nil {

	}

	//execute template with data
	executeTemplate("views/indexWhenAuthorized.html", w, struct {
		User      *User
		ErrUser   error
		Course    *[]Course
		ErrCourse error
		Documents []Document
		ErrDocs   []error
	}{
		User:      user,
		ErrUser:   err2,
		Course:    course,
		ErrCourse: err,
		Documents: documents,
		ErrDocs:   nil,
	})
}

func getCourses() (*[]Course, error) {

	url := "http://www.nbrb.by/api/exrates/rates?periodicity=0"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var course []Course
	err = json.Unmarshal(body, &course)
	return &course, err
}

func executeTemplate(page string, w http.ResponseWriter, data interface{}) {
	tmpl := template.Must(template.ParseFiles(page))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Something happened when the template [page - %v] executed: %v", page, err)
	}
}

func (h *indexHandler) getDocumentsByUser(user *User) ([]Document, error) {
	var docs []Document
	//var errors []error

	for d := range user.Documents {
		id, err := doPrettyId(fmt.Sprint(user.Documents[d]))
		if err != nil {
			return docs, fmt.Errorf("Can't do id for search in database %v: %v ", id, err)
			//errors = append(errors, err)
			//continue
		}
		var elem Document
		err = h.client.findOneById(docColName, id, &elem)
		if err != nil {
			return docs, fmt.Errorf("Can't find document with id %v: %v ", id, err)
			//errors = append(errors, err)
			//continue
		}
		elem.Id = fmt.Sprint(user.Documents[d])
		docs = append(docs, elem)
	}
	return docs, nil
}

func doPrettyId(stringId string) (primitive.ObjectID, error) {
	stringId = stringId[10 : len(stringId)-2]
	return primitive.ObjectIDFromHex(stringId)
}
