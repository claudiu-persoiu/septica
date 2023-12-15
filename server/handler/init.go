package handler

import (
	"fmt"
	"github.com/claudiu-persoiu/septica/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var h = game.NewHub()

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to websockets %v\n", err)
	}

	client := game.NewClient(conn, h)

	go client.Run()
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)
