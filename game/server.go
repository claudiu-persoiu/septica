package game

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

// Server object
type Server struct {
	handler http.Handler
	hub     *Hub
}

var templates = &template.Template{}

func parseTemplates(templateBox embed.FS) error {

	return fs.WalkDir(templateBox, "template", func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() {
			f, _ := templateBox.Open(path)
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(path)
			_, err = templates.New(path).Parse(string(b))
			return err
		}
		return nil
	})
}

// NewServer generate new game server
func NewServer(publicBox, templateBox embed.FS) (*Server, error) {
	server := &Server{hub: NewHub()}

	err := parseTemplates(templateBox)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/", server.pageHandler)
	router.HandleFunc("/simulator", server.simulatorHandler)
	router.Handle("/public/", http.FileServer(http.FS(publicBox)))
	router.HandleFunc("/ws", server.webSocket)

	server.handler = router

	return server, nil
}

// GetHandler returns the server handler
func (p *Server) GetHandler() http.Handler {
	return p.handler
}

func (p *Server) pageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "template/index.html", nil)
}

func (p *Server) simulatorHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "template/simulator.html", nil)
}

func (p *Server) webSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	ws := newClient(w, r, p.hub)

	go ws.waitForMsg()
	go ws.sendMessage()
}
