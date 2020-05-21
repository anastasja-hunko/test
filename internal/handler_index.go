package internal

import (
	"encoding/json"
	"github.com/anastasja-hunko/test/internal/model"
	"io/ioutil"
	"net/http"
	"time"
)

type Course struct {
	Abbreviation string  `json:"Cur_Abbreviation"`
	Rate         float64 `json:"Cur_OfficialRate"`
}

type indexHandler struct {
	serv  *Server
	title string
	page  string
}

func NewIndexHandler(server *Server) *indexHandler {
	return &indexHandler{serv: server, title: "Index"}
}

func (h *indexHandler) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.serv.getUserFromSession(r)

		if err != nil {
			h.serv.Logger.Error(err)
			h.page = "views/indexWhenNonAuthorized.html"
			h.serv.Respond(w, r, http.StatusBadRequest, err, h.page)
			return
		}
		h.page = "views/indexWhenAuthorized.html"
		//get courser from nbrb
		course, err := getCourses()

		if err != nil {
			h.serv.Logger.Error(err)
			h.serv.Respond(w, r, http.StatusBadRequest, err, h.page)
			return
		}

		documents, err := h.serv.DB.Document().FindDocumentsByUser(user)

		if err != nil {
			h.serv.Logger.Error(err)
			h.serv.Respond(w, r, http.StatusBadRequest, err, h.page)
			return
		}

		h.serv.Respond(w, r, http.StatusOK, struct {
			User      model.User
			Course    []Course
			Documents []model.Document
		}{
			User:      *user,
			Course:    *course,
			Documents: *documents,
		}, h.page)
	}
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

func (h *indexHandler) Logout() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.serv.workWithSession(rw, r, "")
		if err != nil {
			h.serv.Logger.Error(err)
			h.serv.Respond(rw, r, http.StatusBadRequest, err, h.page)
			return
		}
		h.serv.Logger.Info("Logout")
		http.Redirect(rw, r, "/", 302)
	}
}
