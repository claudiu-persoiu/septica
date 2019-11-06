package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type game struct {
	State     int
	key       string
	Clients   []*Client
	Deck      []*card
	table     []*card
	firstCard int
}

const (
	WAITING = 0
	STARTED = 1
	OVER    = 2
	BROKEN  = 3
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

	g.Clients = append(g.Clients, client)
	client.position = len(g.Clients) - 1
	client.game = g

	g.notifyClients(&message{Action: "joined", Data: strconv.Itoa(len(g.Clients))})

	return nil
}

func (g *game) notifyClients(m *message) {
	for _, client := range g.Clients {
		client.Send <- m
	}
}

func (g *game) Start(client *Client) error {

	if g.State != WAITING && g.State != OVER {
		return errors.New("started")
	}

	if g.Clients[0] != client {
		return errors.New("not host")
	}

	if len(g.Clients) < 2 {
		return errors.New("not enough players")
	}

	g.State = STARTED
	g.table = nil
	// client that will server
	g.firstCard = 0

	fmt.Println("started game")

	g.Deck = cardsShuffle()

	for _, client := range g.Clients {
		client.cards = append(client.cards, g.Deck[len(g.Deck)-4:]...)
		client.points = 0
		g.Deck = g.Deck[:len(g.Deck)-4]
		cards, _ := json.Marshal(client.cards)
		client.Send <- &message{Action: "cards", Data: string(cards)}
		client.Send <- &message{Action: "possition", Data: strconv.Itoa(client.position)}
		client.Send <- &message{Action: "first", Data: strconv.Itoa(g.firstCard)}
	}

	g.notifyClientsTableUpdate()

	return nil
}

func (g *game) play(client *Client, cardIndex int) error {

	if err := g.validTurn(client); err != nil {
		return err
	}

	if (len(client.cards) - 1) < cardIndex {
		return errors.New("card unavailable")
	}

	card := client.cards[cardIndex]

	tableLen := len(g.table)
	if tableLen > 0 && tableLen%len(g.Clients) == 0 {
		if !g.isCut(card) {
			return errors.New("this card is not a valid cut")
		}
	}

	g.table = append(g.table, card)
	client.cards = append(client.cards[:cardIndex], client.cards[cardIndex+1:]...)

	cards, _ := json.Marshal(client.cards)
	client.Send <- &message{Action: "cards", Data: string(cards)}

	g.notifyClientsTableUpdate()

	return nil
}

func (g *game) isCut(card *card) bool {
	if card.Number == "7" {
		return true
	}

	if len(g.Clients)%3 == 0 && card.Number == "8" {
		return true
	}

	if g.table[0].Number == card.Number {
		return true
	}

	return false
}

func (g *game) validTurn(client *Client) error {
	// see if it's this client's turn
	if (len(g.table)+g.firstCard)%len(g.Clients) != client.position {
		return errors.New("invalid turn")
	}

	return nil
}

func (g *game) getLastPlayerCut() *Client {
	for i := len(g.table) - 1; i >= 0; i-- {
		if g.isCut(g.table[i]) {
			index := (i + g.firstCard) % len(g.Clients)
			return g.Clients[index]
		}
	}

	return g.Clients[0]
}

func (g *game) notifyClientsTableUpdate() {
	cards, _ := json.Marshal(g.table)
	g.notifyClients(&message{Action: "table", Data: string(cards)})
}

func (g *game) fetchHand(client *Client) error {

	if err := g.validTurn(client); err != nil {
		return err
	}

	if len(g.table)%len(g.Clients) != 0 {
		return errors.New("invalid fetch")
	}

	c := g.getLastPlayerCut()
	for _, card := range g.table {
		if card.Number == "10" || card.Number == "A" {
			c.points++
		}
	}

	g.firstCard = c.position

	// publish points?
	g.table = []*card{}
	g.notifyClientsTableUpdate()
	g.notifyClients(&message{Action: "first", Data: strconv.Itoa(c.position)})

	if len(g.Deck) == 0 {
		// no more cards to deal
		if len(g.Clients[0].cards) == 0 {
			// the game is over
			return g.finishGame()
		}

		return nil
	}

	cardsMissing := (4 - len(g.Clients[0].cards)) * len(g.Clients)
	cardsPerPlayer := cardsMissing / len(g.Clients)

	if cardsMissing > len(g.Deck) {
		cardsPerPlayer = len(g.Deck) / len(g.Clients)
	}

	for _, client := range g.Clients {
		client.cards = append(client.cards, g.Deck[len(g.Deck)-cardsPerPlayer:]...)
		g.Deck = g.Deck[:len(g.Deck)-cardsPerPlayer]
		cards, _ := json.Marshal(client.cards)
		client.Send <- &message{Action: "cards", Data: string(cards)}
	}

	return nil
}

func (g *game) finishGame() error {
	g.State = OVER

	g.notifyClients(g.getResultMessage())
	return nil
}

func (g *game) getResultMessage() *message {
	result := make(map[int]int)

	for _, client := range g.Clients {
		result[client.position] = client.points
	}

	resultsString, _ := json.Marshal(result)
	return &message{Action: "result", Data: string(resultsString)}
}

func (g *game) leave() error {
	g.notifyClients(&message{Action: "left"})

	for _, client := range g.Clients {
		client.game = nil
		client.points = 0
		client.cards = nil
	}

	g.State = BROKEN

	return nil
}

func (g *game) restart(client *Client) error {
	g.notifyClients(&message{Action: "restarting"})
	return g.Start(client)
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
