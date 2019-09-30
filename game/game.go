package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	State   int
	Clients []*Client
	Deck    []*Card
	table   []*Card
	turn    int
}

const (
	WAITING = 0
	STARTED = 1
	OVER    = 2
)

func NewGame() *Game {
	g := &Game{State: WAITING}

	return g
}

func (g *Game) AddPlayer(client *Client) error {
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

func (g *Game) notifyClients() {
	for _, client := range g.Clients {
		client.Send <- &Message{Action: "joined"}
	}
}

func (g *Game) Start() {
	g.State = STARTED

	fmt.Println("started game")

	g.Deck = cardsShuffle()

	// cards, _ := json.Marshal(g.Deck)
	// fmt.Println(string(cards))

	for _, client := range g.Clients {
		client.cards = append(client.cards, g.Deck[len(g.Deck)-4:]...)
		g.Deck = g.Deck[:len(g.Deck)-4]
		cards, _ := json.Marshal(client.cards)
		client.Send <- &Message{Action: "cards", Data: string(cards)}
	}

	fmt.Println(len(g.Deck))
}

var deck = []*Card{
	&Card{Number: "7", Type: "diamond"},
	&Card{Number: "8", Type: "diamond"},
	&Card{Number: "9", Type: "diamond"},
	&Card{Number: "10", Type: "diamond"},
	&Card{Number: "J", Type: "diamond"},
	&Card{Number: "Q", Type: "diamond"},
	&Card{Number: "K", Type: "diamond"},
	&Card{Number: "A", Type: "diamond"},
	&Card{Number: "7", Type: "hearts"},
	&Card{Number: "8", Type: "hearts"},
	&Card{Number: "9", Type: "hearts"},
	&Card{Number: "10", Type: "hearts"},
	&Card{Number: "J", Type: "hearts"},
	&Card{Number: "Q", Type: "hearts"},
	&Card{Number: "K", Type: "hearts"},
	&Card{Number: "A", Type: "hearts"},
	&Card{Number: "7", Type: "spades"},
	&Card{Number: "8", Type: "spades"},
	&Card{Number: "9", Type: "spades"},
	&Card{Number: "10", Type: "spades"},
	&Card{Number: "J", Type: "spades"},
	&Card{Number: "Q", Type: "spades"},
	&Card{Number: "K", Type: "spades"},
	&Card{Number: "A", Type: "spades"},
	&Card{Number: "7", Type: "clubs"},
	&Card{Number: "8", Type: "clubs"},
	&Card{Number: "9", Type: "clubs"},
	&Card{Number: "10", Type: "clubs"},
	&Card{Number: "J", Type: "clubs"},
	&Card{Number: "Q", Type: "clubs"},
	&Card{Number: "K", Type: "clubs"},
	&Card{Number: "A", Type: "clubs"},
}

func cardsShuffle() []*Card {
	newDeck := make([]*Card, len(deck))
	copy(newDeck, deck)

	fmt.Println(deck)
	fmt.Println(newDeck)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(newDeck), func(i, j int) { newDeck[i], newDeck[j] = newDeck[j], newDeck[i] })

	return newDeck
}
