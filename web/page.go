package web

import (
	"net/http"
	"log"
	"html/template"
)

type Page struct {
	Title   string
	Address string
	File    string
}

func (page *Page) Handle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("public/template/" + page.File + ".html")

	if err != nil {
		log.Panic(err)
	}

	err = t.Execute(w, *page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}