package internal

import (
	"fmt"
	"github.com/anastasja-hunko/test/internal/model"
	"github.com/gorilla/sessions"

	//"fmt"
	db "github.com/anastasja-hunko/test/internal/database"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

const (
	sessionName     = "session-name"
	sessionLoginKey = "login"
)

type Server struct {
	config       *Config
	Logger       *logrus.Logger
	router       *mux.Router
	DB           *db.Database
	sessionStore sessions.Store
}

func New(config *Config, sessionStore sessions.Store) *Server {
	return &Server{
		config:       config,
		Logger:       logrus.New(),
		router:       mux.NewRouter(),
		sessionStore: sessionStore,
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureDatabase(); err != nil {
		return err
	}

	defer func() {
		if err := s.DB.Close(); err != nil {
			s.Logger.Error("can't close db connection...")
		}
	}()

	s.Logger.Info("starting server...")
	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}

	s.Logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static", http.StripPrefix("/static/", fs))

	autorHandler := NewAutorHandler(s)
	s.router.HandleFunc("/authorization", autorHandler.HandleAuthorize())

	docHandler := NewDocHandler(s)
	s.router.HandleFunc("/createDoc", docHandler.CreateDocHandler())
	s.router.HandleFunc("/editDoc", docHandler.EditDocHandler())
	s.router.HandleFunc("/deleteDoc", docHandler.DeleteDocument())

	indexHandler := NewIndexHandler(s)
	s.router.HandleFunc("/", indexHandler.HandleIndex())
	s.router.HandleFunc("/logout", indexHandler.Logout())

	registerHandler := NewRegHandler(s)
	s.router.HandleFunc("/register", registerHandler.HandleRegister())
}

func (s *Server) configureDatabase() error {
	dbase := db.New(s.config.DbConfig)

	if err := dbase.Open(); err != nil {
		return err
	}

	s.DB = dbase
	return nil
}

func executeTemplate(page string, w http.ResponseWriter, data interface{}) error {
	tmpl := template.Must(template.ParseFiles(page))
	return tmpl.Execute(w, data)
}

func (s *Server) Respond(rw http.ResponseWriter, code int, data interface{}, page string) {
	rw.WriteHeader(code)
	if data != nil {
		if err := executeTemplate(page, rw, data); err != nil {
			s.Logger.Error(err)
			fmt.Fprint(rw, data)
		}
	}
}

func (s *Server) Error(rw http.ResponseWriter, code int, err error) {
	s.Respond(rw, code, err, "views/error.html")
}

func (s *Server) workWithSession(w http.ResponseWriter, r *http.Request, login string) error {
	session, err := s.sessionStore.Get(r, sessionName)

	if err != nil {
		return fmt.Errorf("can't get session with name %v", sessionName)
	}

	session.Values[sessionLoginKey] = login
	return sessions.Save(r, w)
}

func (s *Server) getUserFromSession(r *http.Request) (*model.User, error) {
	session, err := s.sessionStore.Get(r, sessionName)

	if err != nil {
		return nil, fmt.Errorf("can't get session with name %v", sessionName)
	}

	login := session.Values[sessionLoginKey]
	return s.DB.User().FindByLogin(fmt.Sprint(login))
}
