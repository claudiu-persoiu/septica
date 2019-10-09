package game

import (
	"errors"
	"fmt"
	"math/rand"
)

// Hub games hub
type Hub struct {
	games    map[string]*game
	Messages chan string
	users    map[*Client]string
}

// NewHub create new Hub
func NewHub() *Hub {
	hub := &Hub{
		Messages: make(chan string),
		games:    make(map[string]*game),
		users:    make(map[*Client]string)}
	return hub
}

// Start start a new game
func (h *Hub) Start(client *Client) {
	key, ok := h.users[client]
	if !ok {
		g := newGame()
		key = h.registerGame(g)
		g.AddPlayer(client)
		h.users[client] = key
	}
	fmt.Println("Starting game: " + key)
	client.Send <- &message{Action: "start", Data: key}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *Hub) registerGame(game *game) string {
	key := randSeq(7)
	h.games[key] = game

	return key
}

func (h *Hub) join(gameKey string, client *Client) error {
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

func (h *Hub) begin(client *Client) error {
	g, err := getGameFromClient(h, client)

	if err != nil {
		return errors.New("game not found")
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

func (h *Hub) play(client *Client, cardIndex int) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return errors.New("game not found")
	}

	// see if it's this client's turn
	if len(g.table)%len(g.Clients) != client.position {
		return errors.New("invalid turn")
	}

	card := client.cards[cardIndex]

	// see if the card is available to the client
	if card == nil {
		return errors.New("card unavailable")
	}

	g.table = append(g.table, card)
	client.cards = append(client.cards[:cardIndex], client.cards[cardIndex+1:]...)

	// publish table to all users

	return nil
}

func getGameFromClient(h *Hub, c *Client) (*game, error) {
	gKey, ok := h.users[c]

	if ok == false {
		return nil, errors.New("invalid")
	}

	g, ok := h.games[gKey]

	if ok == false {
		return nil, errors.New("invalid")
	}

	return g, nil
}
