package internal

import (
	"github.com/anastasja-hunko/test/internal/model"
	"github.com/pkg/errors"
	"net/http"
)

type regHandler struct {
	serv  *Server
	title string
	page  string
}

func NewRegHandler(serv *Server) *regHandler {
	return &regHandler{serv: serv, title: "Registration", page: "views/userForm.html"}
}

func (h *regHandler) HandleRegister() http.HandlerFunc {
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

			if err := h.registerUser(u); err != nil {
				template := template{
					Pagetitle: h.title,
					Message:   "Something went wrong!",
					Err:       err,
				}
				h.serv.Logger.Error(err)
				h.serv.Respond(rw, r, http.StatusBadRequest, template, h.page)
				return
			}
			template := template{
				Pagetitle: h.title,
				Message:   "User was created!",
				Err:       nil,
			}
			h.serv.Logger.Info("User was created...")
			h.serv.Respond(rw, r, http.StatusCreated, template, h.page)
		}
		template := template{
			Pagetitle: h.title,
		}
		h.serv.Respond(rw, r, http.StatusOK, template, h.page)
	}
}

func (h *regHandler) registerUser(u *model.User) error {
	user, err := h.serv.DB.User().FindByLogin(u.Login)

	if user != nil {
		return errors.New("user's already existed with login:" + u.Login)
	}

	err = h.serv.DB.User().Create(u)

	if err != nil {
		return err
	}

	return nil
}
