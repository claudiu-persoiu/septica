package game

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	SetWriteDeadline(t time.Time) error
	WriteMessage(messageType int, data []byte) error
	NextWriter(messageType int) (io.WriteCloser, error)
	Close() error
}

// Client - game Client
type Client struct {
	connection Connection
	name       string
	identifier string
	send       chan *message
	hub        *hub
	game       *game
	cards      []*card
	position   int
	points     int
	won        int
}

func NewClient(conn Connection, hub *hub) *Client {
	return &Client{connection: conn, send: make(chan *message, 256), hub: hub}
}

func (c *Client) processMessage(m message) {

	switch m.Action {
	case "name":
		c.name = m.Data
		c.send <- &message{Action: "name", Data: c.name}
		c.send <- &message{Action: "nogame"}
	case "identify":
		identifier := m.Data
		client, _ := c.hub.users[identifier]

		c.hub.users[identifier] = c
		c.identifier = identifier

		if client != nil {
			c.game = client.game
			c.cards = client.cards
			c.position = client.position
			c.points = client.points
			c.name = client.name
		}

		if c.name != "" {
			c.send <- &message{Action: "name", Data: c.name}
		} else {
			c.send <- &message{Action: "noname"}
			return
		}

		if c.game != nil {
			c.game.Clients[c.position] = c

			c.send <- &message{Action: "position", Data: strconv.Itoa(c.position)}
			c.send <- &message{Action: "joined", Data: strconv.Itoa(len(c.game.Clients))}
			c.send <- &message{Action: "start", Data: c.game.key}

			namesJSON, _ := json.Marshal(c.game.GetNames())
			c.send <- &message{Action: "names", Data: string(namesJSON)}

			if c.game.State == STARTED {
				c.send <- &message{Action: "first", Data: strconv.Itoa(c.game.firstCard)}

				cards, _ := json.Marshal(c.game.table)
				c.send <- &message{Action: "table", Data: string(cards)}

				cards, _ = json.Marshal(c.cards)
				c.send <- &message{Action: "cards", Data: string(cards)}
			} else if c.game.State == OVER {
				cards, _ := json.Marshal(c.game.table)
				c.send <- &message{Action: "table", Data: string(cards)}

				c.send <- c.game.getResultMessage()
				c.send <- c.game.getGamesStatsMessage()
			}
		} else {
			c.send <- &message{Action: "nogame"}
		}

	case "start":
		c.hub.Start(c)
	case "join":
		err := c.hub.join(m.Data, c)
		if err != nil {
			c.send <- &message{Action: "join", Data: err.Error()}
		}
	case "begin":
		err := c.hub.begin(c)
		if err != nil {
			c.send <- &message{Action: "error", Data: err.Error()}
		}
	case "play":
		i, err := strconv.Atoi(m.Data)
		if err != nil {
			c.send <- &message{Action: "error", Data: "invalid card index send"}
		} else {
			err := c.hub.play(c, i)
			if err != nil {
				c.send <- &message{Action: "error", Data: err.Error()}
			}
		}
	case "fetch":
		err := c.hub.fetchHand(c)
		if err != nil {
			c.send <- &message{Action: "error", Data: err.Error()}
		}
	case "leave":
		c.hub.leave(c)
	case "restart":
		c.hub.restartGame(c)
	default:
		c.send <- &message{Action: "error", Data: "invalid command"}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func (c *Client) Run() {
	go c.waitForMsg()
	go c.sendMessage()
}

func (c *Client) Send(message *message) {
	c.send <- message
}

func (c *Client) waitForMsg() {
	defer func() {
		c.connection.Close()
		// TODO handle user disconnect to stop game
		fmt.Println("User disconnected")
	}()
	for {
		_, msg, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Print(string(msg))

		var obj message
		if err := json.Unmarshal(msg, &obj); err == nil {
			c.processMessage(obj)
			log.Print(obj.Action)
		} else {
			log.Println("Error parsing message:")
			log.Println(err)
		}
	}
}

func (c *Client) sendMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// @TODO handle user disconnect to stop game
		ticker.Stop()
		c.connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:

			fmt.Println(message)

			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			jsonMessage, err := json.Marshal(message)
			if err != nil {
				return
			}

			w.Write(jsonMessage)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
