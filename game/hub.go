package game

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

// Hub games hub
type Hub struct {
	games    map[string]*game
	Messages chan string
	users    map[string]*Client
}

// NewHub create new Hub
func NewHub() *Hub {
	hub := &Hub{
		Messages: make(chan string),
		games:    make(map[string]*game),
		users:    make(map[string]*Client)}
	return hub
}

// Start start a new game
func (h *Hub) Start(client *Client) error {
	key := ""
	if client.game == nil {
		g := newGame()
		key = h.registerGame(g)
		if err := h.join(key, client); err != nil {
			return err
		}
	} else if client.game.State == WAITING || client.game.State == STARTED {
		key = client.game.key
	}
	fmt.Println("Starting game: " + key)
	client.Send <- &message{Action: "start", Data: key}

	return nil
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
	game.key = key

	return key
}

func (h *Hub) join(gameKey string, client *Client) error {
	g, ok := h.games[gameKey]

	if ok == false {
		return errors.New("invalid")
	}

	err := g.AddPlayer(client)

	if err == nil {
		client.game = g
	}

	client.Send <- &message{Action: "start", Data: g.key}
	client.Send <- &message{Action: "position", Data: strconv.Itoa(client.position)}

	return err
}

func (h *Hub) begin(client *Client) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.Start(client, 0)
}

func (h *Hub) play(client *Client, cardIndex int) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.play(client, cardIndex)
}

func (h *Hub) fetchHand(client *Client) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.fetchHand(client)
}

func (h *Hub) leave(client *Client) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	err = g.leave()

	if err != nil {
		return err
	}

	delete(h.games, g.key)

	return nil
}

func (h *Hub) restartGame(client *Client) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.restart(client)
}

func getGameFromClient(h *Hub, c *Client) (*game, error) {
	if c.identifier == "" {
		return nil, errors.New("unidentified user")
	}

	_, ok := h.users[c.identifier]

	if ok == false {
		return nil, errors.New("invalid user in hub")
	}

	if c.game == nil {
		return nil, errors.New("invalid game")
	}

	return c.game, nil
}
