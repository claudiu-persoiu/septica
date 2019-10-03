package game

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Server object
type Server struct {
	handler http.Handler
	hub     *Hub
}

// NewServer generate new game server
func NewServer() (*Server, error) {
	server := &Server{hub: NewHub()}

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.pageHandler))
	router.Handle("/simulator", http.HandlerFunc(server.simulatorHandler))
	router.Handle("/js/", http.FileServer(http.Dir("public")))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))

	server.handler = router

	return server, nil
}

// GetHandler returns the server handler
func (p *Server) GetHandler() http.Handler {
	return p.handler
}

func (p *Server) pageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/template/index.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	t.Execute(w, nil)
}

func (p *Server) simulatorHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/template/simulator.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	t.Execute(w, nil)
}

func (p *Server) webSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	ws := newClient(w, r, p.hub)

	go ws.waitForMsg()
	go ws.sendMessage()
}
