package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/claudiu-persoiu/sedma/server"
)

func main() {

	address := ":8080"

	fmt.Println("Starting server: " + address)
	server, err := server.NewGameServer()

	if err != nil {
		log.Fatal("Unable to start server", err)
	}

	err = http.ListenAndServe(address, server)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
