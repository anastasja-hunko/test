package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Course struct {
	Abbreviation string  `json:"Cur_Abbreviation"`
	Rate         float64 `json:"Cur_OfficialRate"`
}

type socketHandler struct {
	serv *Server
}

type wss struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

func NewSocketHandler(serv *Server) *socketHandler {
	return &socketHandler{serv: serv}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *socketHandler) Try() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Client subscribed")

		wss := wss{Socket: conn}

		exchange := []string{"USD", "EUR", "JPY"}

		for {
			for _, r := range exchange {
				time.Sleep(2 * time.Second)
				go h.getCourse(r, &wss)
			}
		}
	}
}

func (h *socketHandler) getCourse(r string, conn *wss) {

	time.Sleep(2 * time.Second)

	url := "http://www.nbrb.by/api/exrates/rates/" + r + "?parammode=2"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		h.serv.Logger.Error(err.Error())
	}

	res, err := client.Do(req)

	if err != nil {
		h.serv.Logger.Error(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		h.serv.Logger.Error(err.Error())
	}

	var c Course

	err = json.Unmarshal(body, &c)

	if err != nil {
		h.serv.Logger.Error(err.Error())
	}

	err = conn.send([]byte((fmt.Sprint(c.Abbreviation, " : ", c.Rate))))

	if err != nil {
		h.serv.Logger.Error(err.Error())
	}
}

func (p *wss) send(v []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Socket.WriteMessage(websocket.TextMessage, v)
}
