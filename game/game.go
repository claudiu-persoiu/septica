package game

type Game struct {
	State int
	Clients []*Client
	Deck []Card
	turn int
}

const WAITING = 0
const START = 1


func NewGame() *Game {
	g := &Game{State: WAITING}

	return g
}

func (g *Game) AddPlayer(client *Client)  {
	if len(g.Clients) < 4 {
		g.Clients = append(g.Clients, client)
	}
}