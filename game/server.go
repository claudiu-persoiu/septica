package game

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type GameServer struct {
	handler http.Handler
	hub     *Hub
}

// NewGameServer generate new game server
func NewGameServer() (*GameServer, error) {
	server := &GameServer{hub: NewHub()}

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.pageHandler))
	router.Handle("/simulator", http.HandlerFunc(server.simulatorHandler))
	router.Handle("/js/", http.FileServer(http.Dir("public")))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))

	server.handler = router

	return server, nil
}

// GetHandler returns the server handler
func (p *GameServer) GetHandler() http.Handler {
	return p.handler
}

func (p *GameServer) pageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/template/index.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	t.Execute(w, nil)
}

func (p *GameServer) simulatorHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/template/simulator.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	t.Execute(w, nil)
}

func (p *GameServer) webSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	ws := newClient(w, r, p.hub)

	go ws.waitForMsg()
	go ws.sendMessage()
}
