package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type playerServerWS struct {
	*websocket.Conn
}

type message struct {
	Action int    `json:"action"`
	Data   string `json:"data,omitempty"`
}

func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(1, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to websockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}

	log.Print(string(msg))

	var obj message
	if err := json.Unmarshal(msg, &obj); err == nil {
		// c.processMessage(obj)
		log.Print(obj.Action)
	} else {
		log.Fatal(err)
	}
}
