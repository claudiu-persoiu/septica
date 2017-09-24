package game

import (
	"math/rand"
)

type Hub struct {
	games map[string]*Game
	Messages chan string
	Start chan *Client
	users map[*Client]string
}

func NewHub() *Hub {
	hub := &Hub{
		Messages: make(chan string),
		Start: make(chan *Client),
		games: make(map[string] *Game),
		users: make(map[*Client]string)}
	go hub.processMessages()
	return hub
}

func (h *Hub) processMessages() {
	for {
		select {
		case client := <- h.Start:
			key, ok := h.users[client]
			if !ok {
				g := NewGame()
				key = h.registerGame(g)
				g.AddPlayer(client)
				h.users[client] = key
			}
			client.Send <- &Message{Step:START, Command:key}
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
	key := randSeq(7)
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