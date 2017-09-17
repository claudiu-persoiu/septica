package web

import (
	"net/http"
	"log"
	"html/template"
)

type Page struct {
	Title   string
	Address string
}

func (page *Page) Handle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/template/index.html")

	if err != nil {
		log.Panic(err)
	}

	err = t.Execute(w, *page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
