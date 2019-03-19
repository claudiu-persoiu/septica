package main

import (
	"log"
	"net/http"

	"github.com/claudiu-persoiu/sedma/server"
)

func main() {

	// websocketPath := "/echo"
	address := ":8080"

	// page := web.NewPageHandler("Sedman", address+websocketPath, "index")
	// http.HandleFunc("/", page.Handle)

	// page = web.NewPageHandler("Sedman Simulator", address+websocketPath, "simulator")
	// http.HandleFunc("/simulator", page.Handle)

	// http.Handle("/js/", http.FileServer(http.Dir("public")))

	// hub := game.NewHub()

	// http.HandleFunc(websocketPath, func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("joined...")
	// 	game.HandleWebsocket(w, r, hub)
	// })

	// fmt.Println("Starting server: " + address)

	// server, err := NewGameServer()

	server, err := server.NewGameServer()

	if err != nil {
		log.Fatal("Unable to start server", err)
	}

	err = http.ListenAndServe(address, server)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
