package game

import (
	"math/rand"
)

type Hub struct {
	games map[string]*Game
	Messages chan string
	// connected users that are not playing yet
	lobby map[*Client]bool
	Start chan *Client
}

func NewHub() *Hub {
	hub := &Hub{Messages: make(chan string), lobby: make(map[*Client]bool)}
	go hub.processMessages()
	return hub
}

func (h *Hub) processMessages() {
	for {
		select {
		case client := <- h.Start:
			g := NewGame()
			key := h.registerGame(g)
			g.AddPlayer(client)
			client.send <- []byte(key)
		}
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *Hub) registerGame(game *Game) string {
	key := randSeq(5)
	h.games[key] = game

	return key
}

// maybe refactor into a channel?
func (h *Hub) Join(gameKey string) *Game {
	g, ok := h.games[gameKey]

	if ok == false {
		return nil
	}

	return g
}