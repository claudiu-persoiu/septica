package main

import (
	"embed"
	"fmt"
	"github.com/claudiu-persoiu/septica/server/router"
	"github.com/claudiu-persoiu/septica/server/templates"
	"log"
	"net/http"
)

//go:embed public
var publicBox embed.FS

//go:embed template
var templateBox embed.FS

func main() {
	address := ":8008"

	fmt.Println("Starting server: " + address)

	templatesHandler := templates.NewTemplatesHandler(templateBox)
	publicHandler := http.FileServer(http.FS(publicBox))

	r := router.Init(templatesHandler, publicHandler.ServeHTTP)

	err := http.ListenAndServe(address, r)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
