package main

import (
	"net/http"
	"log"
	"github.com/claudiu-persoiu/sedma/web"
	"fmt"
	"github.com/claudiu-persoiu/sedma/game"
)

func main() {

	websocketPath := "/echo"
	address := ":8000"

	page := &web.Page{Title: "Sedman", Address: address + websocketPath}

	http.HandleFunc("/", page.Handle)
	http.Handle("/js/", http.FileServer(http.Dir("public")))

	hub := game.NewHub()

	http.HandleFunc(websocketPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("joined...")
		game.HandleWebsocket(w, r, hub)
	})

	fmt.Println("Starting server: " + address)

	err := http.ListenAndServe(address, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
