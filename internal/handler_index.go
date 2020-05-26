package internal

import (
	"github.com/anastasja-hunko/test/internal/model"
	"net/http"
)

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
			h.serv.Respond(w, http.StatusOK, err, h.page)
			return
		}
		h.page = "views/indexWhenAuthorized.html"

		documents, err := h.serv.DB.Document().FindDocumentsByUser(user)

		if err != nil {
			h.serv.Logger.Error(err)
			h.serv.Error(w, http.StatusBadRequest, err)
			return
		}

		h.serv.Respond(w, http.StatusOK, struct {
			User      model.User
			Documents []model.Document
		}{
			User:      *user,
			Documents: *documents,
		}, h.page)
	}
}

func (h *indexHandler) Logout() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.serv.workWithSession(rw, r, "")
		if err != nil {
			h.serv.Logger.Error(err)
			h.serv.Error(rw, http.StatusBadRequest, err)
			return
		}
		h.serv.Logger.Info("Logout")
		http.Redirect(rw, r, "/", 302)
	}
}
