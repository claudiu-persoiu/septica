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
	router := game.InitRouter(publicBox, templateBox)

	err := http.ListenAndServe(address, router)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
