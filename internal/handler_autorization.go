package internal

import (
	"errors"
	"github.com/anastasja-hunko/test/internal/model"
	"net/http"
)

type autorHandler struct {
	serv  *Server
	title string
	page  string
}

func NewAutorHandler(serv *Server) *autorHandler {
	return &autorHandler{serv: serv, title: "Authorization", page: "views/userForm.html"}
}

func (h *autorHandler) HandleAuthorize() http.HandlerFunc {
	type template struct {
		Pagetitle string
		Message   string
		Err       error
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			u := &model.User{
				Login:    r.FormValue("login"),
				Password: r.FormValue("password"),
			}

			if err := h.authorize(u); err != nil {
				template := template{
					Pagetitle: h.title,
					Message:   "Something went wrong!",
					Err:       err,
				}
				h.serv.Logger.Error("User was not authorized: ", err)
				h.serv.Respond(rw, http.StatusUnauthorized, template, h.page)
				return
			}
			h.serv.Logger.Info("User was authorized... ")

			err := h.serv.workWithSession(rw, r, u.Login)

			if err != nil {
				h.serv.Logger.Error("Problem with session: ", err)
			}
			h.serv.Logger.Info("User login was stored in the session: ")

			http.Redirect(rw, r, "/", 302)
			return
		}
		template := template{
			Pagetitle: h.title,
		}
		h.serv.Respond(rw, http.StatusOK, template, h.page)
	}
}

func (h *autorHandler) authorize(u *model.User) error {
	user, err := h.serv.DB.User().FindByLogin(u.Login)

	if err != nil || !user.ComparePasswords(u.Password) {
		return errors.New("incorrect password or login")
	}

	return nil
}
