package internal

import (
	"fmt"
	"github.com/anastasja-hunko/test/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
)

type docHandler struct {
	server *Server
	user   *model.User
	page   string
}

func NewDocHandler(s *Server) *docHandler {
	return &docHandler{server: s, page: "views/createEdit.html"}
}

func (h *docHandler) CreateDocHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user, err := h.server.getUserFromSession(r)
		if err != nil {
			h.server.Logger.Error(err)
			h.server.Error(rw, http.StatusBadRequest, err)
			return
		}

		h.user = user

		if r.Method == http.MethodPost {
			doc := model.Document{Title: r.FormValue("Title"), Content: r.FormValue("Content")}
			err := h.create(&doc)
			if err != nil {
				h.server.Logger.Error(err)
				h.server.Error(rw, http.StatusBadRequest, err)
				return
			}
			h.server.Logger.Info("document was created")
			http.Redirect(rw, r, "/", 302)
			return
		}

		h.showDocForm(nil, rw)
	}
}

func (h *docHandler) showDocForm(doc *model.Document, rw http.ResponseWriter) {
	docView := model.DocView{DocumentInput: getDocumentInput(doc), User: h.user}
	h.server.Respond(rw, http.StatusOK, docView, h.page)
}

func (h *docHandler) create(doc *model.Document) error {
	insertedResult, err := h.server.DB.Document().Create(doc)
	if err != nil {
		return fmt.Errorf("can't insert a document: %v", err)
	}

	docs := h.user.Documents
	docs = append(docs, insertedResult.InsertedID)

	return h.server.DB.User().SetUserDocs(h.user, docs)
}

func (h *docHandler) EditDocHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user, err := h.server.getUserFromSession(r)
		if err != nil {
			h.server.Logger.Error(err)
			h.server.Error(rw, http.StatusBadRequest, err)
			return
		}

		h.user = user

		doc, err := h.getDocument(r.URL)

		if err != nil {
			h.server.Logger.Error(err)
			h.server.Error(rw, http.StatusBadRequest, err)
			return
		}

		if r.Method == http.MethodPost {
			doc.Title = r.FormValue("Title")
			doc.Content = r.FormValue("Content")

			err := h.server.DB.Document().Edit(doc)

			if err != nil {
				h.server.Logger.Error(err)
				h.server.Error(rw, http.StatusBadRequest, err)
				return
			}

			h.server.Logger.Info("document was edited")
			http.Redirect(rw, r, "/", 302)
			return
		}

		h.showDocForm(doc, rw)
	}
}

func (h *docHandler) DeleteDocument() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id, err := getDocIdFromRequest(r.URL)

		if err != nil {
			h.server.Logger.Error("can't get if from request: ", err)
			h.server.Error(rw, http.StatusBadRequest, err)
			return
		}
		err = h.server.DB.Document().Delete(id)

		if err != nil {
			h.server.Logger.Error("can't delete a document: ", err)
			h.server.Error(rw, http.StatusBadRequest, err)
			return
		}

		err = h.server.DB.User().RemoveIdFromUserDocs(h.user, id)

		if err != nil {
			h.server.Logger.Error("document was deleted from collection docs, but a connection with user still present: ", err)
			h.server.Error(rw, http.StatusBadRequest, err)
			return
		}

		http.Redirect(rw, r, "/", 302)
	}
}

func (h *docHandler) getDocument(url *url.URL) (*model.Document, error) {
	objectId, err := getDocIdFromRequest(url)

	if err != nil {
		return nil, err
	}

	doc, err := h.server.DB.Document().FindById(objectId)

	if err != nil {
		return nil, err
	}

	doc.Id = objectId
	return doc, nil
}

func getDocIdFromRequest(url *url.URL) (primitive.ObjectID, error) {
	id := fmt.Sprint(url.Query().Get("docId"))
	return doPrettyId(id)
}

func doPrettyId(stringId string) (primitive.ObjectID, error) {
	stringId = stringId[10 : len(stringId)-2]
	return primitive.ObjectIDFromHex(stringId)
}

func getDocumentInput(doc *model.Document) *model.DocumentInput {
	title := ""
	content := ""

	if doc != nil {
		title = doc.Title
		content = doc.Content
	}
	var inputs []model.Input
	input := model.Input{Name: "Title", Caption: "Enter title", Type: "input", Value: title, Required: true}
	inputs = append(inputs, input)
	input2 := model.Input{Name: "Content", Caption: "Enter content", Type: "textarea", Value: content, Required: true}
	inputs = append(inputs, input2)

	return &model.DocumentInput{Inputs: &inputs, Create: doc == nil}
}
