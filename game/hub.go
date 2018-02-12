package game

import (
	"math/rand"
	"errors"
)

type Hub struct {
	games    map[string]*Game
	Messages chan string
	Start    chan *Client
	users    map[*Client]string
}

func NewHub() *Hub {
	hub := &Hub{
		Messages: make(chan string),
		Start:    make(chan *Client),
		games:    make(map[string]*Game),
		users:    make(map[*Client]string)}
	go hub.processMessages()
	return hub
}

func (h *Hub) processMessages() {
	for {
		select {
		case client := <-h.Start:
			key, ok := h.users[client]
			if !ok {
				g := NewGame()
				key = h.registerGame(g)
				g.AddPlayer(client)
				h.users[client] = key
			}
			client.Send <- &Message{Action: "start", Data: key}
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
func (h *Hub) Join(gameKey string, client *Client) error {
	g, ok := h.games[gameKey]

	if ok == false {
		return errors.New("invalid")
	}

	err := g.AddPlayer(client)

	if err != nil {
		h.users[client] = gameKey
	}

	return err
}

func (h *Hub) Begin(client *Client) error {
	gKey, ok := h.users[client]

	if ok == false {
		return errors.New("invalid")
	}

	g, ok := h.games[gKey]

	if ok == false {
		return errors.New("invalid")
	}

	if g.State != WAITING {
		return errors.New("started")
	}

	if g.Clients[0] != client {
		return errors.New("not host")
	}

	g.Start()

	return nil
}