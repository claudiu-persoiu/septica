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

var templates = &template.Template{}
var hub = NewHub()

// InitRouter generate new server router
func InitRouter(publicBox, templateBox embed.FS) http.Handler {
	err := parseTemplates(templateBox)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/", serveTemplate("template/index.html"))
	router.HandleFunc("/simulator", serveTemplate("template/simulator.html"))
	router.Handle("/public/", http.FileServer(http.FS(publicBox)))
	router.HandleFunc("/ws", handleWebSocket)

	return router
}

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

func serveTemplate(template string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, template, nil)
		if err != nil {
			log.Fatal("could not parse template " + template)
			return
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected")
	ws := newClient(w, r, hub)

	go ws.waitForMsg()
	go ws.sendMessage()
}
