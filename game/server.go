package game

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
)

// Server object
type Server struct {
	handler http.Handler
	hub     *Hub
}

var publicBox = packr.New("public", "../public")
var templateBox = packr.New("template", "../template")

var templates = &template.Template{}

func init() {
	templateBox.Walk(func(name string, file file.File) error {
		r := bufio.NewReader(file)
		data, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		_, err = templates.New(name).Parse(string(data))
		return err
	})
}

// NewServer generate new game server
func NewServer() (*Server, error) {
	server := &Server{hub: NewHub()}

	router := http.NewServeMux()
	router.HandleFunc("/", server.pageHandler)
	router.HandleFunc("/simulator", server.simulatorHandler)
	router.HandleFunc("/static/", http.StripPrefix("/static/", http.FileServer(publicBox)).ServeHTTP)
	router.HandleFunc("/ws", server.webSocket)

	server.handler = router

	return server, nil
}

// GetHandler returns the server handler
func (p *Server) GetHandler() http.Handler {
	return p.handler
}

func (p *Server) pageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func (p *Server) simulatorHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "simulator.html", nil)
}

func (p *Server) webSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	ws := newClient(w, r, p.hub)

	go ws.waitForMsg()
	go ws.sendMessage()
}
