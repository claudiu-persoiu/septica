package game

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gobuffalo/packr/v2"
)

// Server object
type Server struct {
	handler http.Handler
	hub     *Hub
}

var publicBox = packr.New("public", "../public")
var templateBox = packr.New("template", "../template")

// NewServer generate new game server
func NewServer() (*Server, error) {
	server := &Server{hub: NewHub()}

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.pageHandler))
	router.Handle("/simulator", http.HandlerFunc(server.simulatorHandler))
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(publicBox)))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))

	server.handler = router

	return server, nil
}

// GetHandler returns the server handler
func (p *Server) GetHandler() http.Handler {
	return p.handler
}

func (p *Server) pageHandler(w http.ResponseWriter, r *http.Request) {
	s, err := templateBox.FindString("index.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	t, err := template.New("hello").Parse(s)

	if err != nil {
		log.Panic(err)
	}

	t.Execute(w, nil)
}

func (p *Server) simulatorHandler(w http.ResponseWriter, r *http.Request) {
	s, err := templateBox.FindString("public/template/simulator.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	t, err := template.New("simulator").Parse(s)

	if err != nil {
		log.Panic(err)
	}

	t.Execute(w, nil)
}

func (p *Server) webSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	ws := newClient(w, r, p.hub)

	go ws.waitForMsg()
	go ws.sendMessage()
}
