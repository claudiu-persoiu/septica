package templates

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

type Handler struct {
	templates *template.Template
}

func NewTemplatesHandler(templateBox embed.FS) *Handler {
	th := &Handler{templates: &template.Template{}}
	err := th.parseTemplates(templateBox)
	if err != nil {
		return nil
	}
	return th
}

func (th *Handler) parseTemplates(templateBox embed.FS) error {
	return fs.WalkDir(templateBox, "template", func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() {
			f, _ := templateBox.Open(path)
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(path)
			_, err = th.templates.New(path).Parse(string(b))
			return err
		}
		return nil
	})
}

func (th *Handler) ServeTemplate(template string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := th.templates.ExecuteTemplate(w, template, data)
		if err != nil {
			log.Fatal("could not parse template " + template)
			return
		}
	}
}
