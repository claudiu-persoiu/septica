package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/claudiu-persoiu/septica/game"
)

//go:embed public
var publicBox embed.FS

//go:embed template
var templateBox embed.FS

func main() {
	address := ":8008"

	fmt.Println("Starting server: " + address)
	server, err := game.NewServer(publicBox, templateBox)

	if err != nil {
		log.Fatal("Unable to start server", err)
	}

	err = http.ListenAndServe(address, server.GetHandler())

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
