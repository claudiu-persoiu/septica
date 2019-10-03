package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/claudiu-persoiu/sedma/game"
)

func main() {

	address := ":8008"

	fmt.Println("Starting server: " + address)
	server, err := game.NewServer()

	if err != nil {
		log.Fatal("Unable to start server", err)
	}

	err = http.ListenAndServe(address, server.GetHandler())

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
