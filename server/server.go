package server

import (
	"log"
	"net/http"
	"text/template"
)

type GameServer struct {
	http.Handler
}

// NewGameServer generate new game server
func NewGameServer() (*GameServer, error) {
	server := new(GameServer)

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.pageHandler))
	router.Handle("/simulator", http.HandlerFunc(server.simulatorHandler))
	// router.Handle("/js/", http.FileServer(http.Dir("public")))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))

	server.Handler = router

	return server, nil
}

func (p *GameServer) pageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/template/index.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	tmpl.Execute(w, nil)
}

func (p *GameServer) simulatorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/template/simulator.html")
	if err != nil {
		log.Fatal("unable to parse template")
	}

	tmpl.Execute(w, nil)
}

func (p *GameServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	ws.WaitForMsg()

	// ws := newPlayerServerWS(w, r)

	// numberOfPlayersMsg := ws.WaitForMsg()
	// numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	// p.game.Start(numberOfPlayers, ws)

	// winner := ws.WaitForMsg()
	// p.game.Finish(winner)
}
