package game

import "errors"

type Game struct {
	State   int
	Clients []*Client
	Deck    []Card
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
	// generat pachet de carti
	// shuffle
	// impartit cartile
	// distribuit carti
}
