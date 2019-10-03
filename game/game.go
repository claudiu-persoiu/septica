package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type game struct {
	State   int
	Clients []*Client
	Deck    []*card
	table   []*card
	turn    int
}

const (
	WAITING = 0
	STARTED = 1
	OVER    = 2
)

func newGame() *game {
	g := &game{State: WAITING}

	return g
}

func (g *game) AddPlayer(client *Client) error {
	if g.State != WAITING {
		return errors.New("started")
	}

	if len(g.Clients) == 4 {
		return errors.New("full")
	}

	g.notifyClients()
	g.Clients = append(g.Clients, client)

	return nil
}

func (g *game) notifyClients() {
	for _, client := range g.Clients {
		client.Send <- &message{Action: "joined"}
	}
}

func (g *game) Start() {
	g.State = STARTED

	fmt.Println("started game")

	g.Deck = cardsShuffle()

	// cards, _ := json.Marshal(g.Deck)
	// fmt.Println(string(cards))

	for _, client := range g.Clients {
		client.cards = append(client.cards, g.Deck[len(g.Deck)-4:]...)
		g.Deck = g.Deck[:len(g.Deck)-4]
		cards, _ := json.Marshal(client.cards)
		client.Send <- &message{Action: "cards", Data: string(cards)}
	}

	fmt.Println(len(g.Deck))
}

var deck = []*card{
	&card{Number: "7", Type: "diamond"},
	&card{Number: "8", Type: "diamond"},
	&card{Number: "9", Type: "diamond"},
	&card{Number: "10", Type: "diamond"},
	&card{Number: "J", Type: "diamond"},
	&card{Number: "Q", Type: "diamond"},
	&card{Number: "K", Type: "diamond"},
	&card{Number: "A", Type: "diamond"},
	&card{Number: "7", Type: "hearts"},
	&card{Number: "8", Type: "hearts"},
	&card{Number: "9", Type: "hearts"},
	&card{Number: "10", Type: "hearts"},
	&card{Number: "J", Type: "hearts"},
	&card{Number: "Q", Type: "hearts"},
	&card{Number: "K", Type: "hearts"},
	&card{Number: "A", Type: "hearts"},
	&card{Number: "7", Type: "spades"},
	&card{Number: "8", Type: "spades"},
	&card{Number: "9", Type: "spades"},
	&card{Number: "10", Type: "spades"},
	&card{Number: "J", Type: "spades"},
	&card{Number: "Q", Type: "spades"},
	&card{Number: "K", Type: "spades"},
	&card{Number: "A", Type: "spades"},
	&card{Number: "7", Type: "clubs"},
	&card{Number: "8", Type: "clubs"},
	&card{Number: "9", Type: "clubs"},
	&card{Number: "10", Type: "clubs"},
	&card{Number: "J", Type: "clubs"},
	&card{Number: "Q", Type: "clubs"},
	&card{Number: "K", Type: "clubs"},
	&card{Number: "A", Type: "clubs"},
}

func cardsShuffle() []*card {
	newDeck := make([]*card, len(deck))
	copy(newDeck, deck)

	fmt.Println(deck)
	fmt.Println(newDeck)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(newDeck), func(i, j int) { newDeck[i], newDeck[j] = newDeck[j], newDeck[i] })

	return newDeck
}
