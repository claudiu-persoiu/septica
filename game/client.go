package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Client - Game client
type Client struct {
	connection *websocket.Conn
	name       string
	identifier string
	Send       chan *message
	hub        *Hub
	game       *game
	cards      []*card
	position   int
	points     int
	won        int
}

func newClient(w http.ResponseWriter, r *http.Request, hub *Hub) *Client {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to websockets %v\n", err)
	}

	return &Client{connection: conn, Send: make(chan *message, 256), hub: hub}
}

func (c *Client) processMessage(m message) {

	switch m.Action {
	case "name":
		c.name = m.Data
		c.Send <- &message{Action: "name", Data: c.name}
		c.Send <- &message{Action: "nogame"}
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
			c.Send <- &message{Action: "name", Data: c.name}
		} else {
			c.Send <- &message{Action: "noname"}
			return
		}

		if c.game != nil {
			c.game.Clients[c.position] = c

			c.Send <- &message{Action: "position", Data: strconv.Itoa(c.position)}
			c.Send <- &message{Action: "joined", Data: strconv.Itoa(len(c.game.Clients))}
			c.Send <- &message{Action: "start", Data: c.game.key}

			namesJSON, _ := json.Marshal(c.game.GetNames())
			c.Send <- &message{Action: "names", Data: string(namesJSON)}

			if c.game.State == STARTED {
				c.Send <- &message{Action: "first", Data: strconv.Itoa(c.game.firstCard)}

				cards, _ := json.Marshal(c.game.table)
				c.Send <- &message{Action: "table", Data: string(cards)}

				cards, _ = json.Marshal(c.cards)
				c.Send <- &message{Action: "cards", Data: string(cards)}
			} else if c.game.State == OVER {
				cards, _ := json.Marshal(c.game.table)
				c.Send <- &message{Action: "table", Data: string(cards)}

				c.Send <- c.game.getResultMessage()
				c.Send <- c.game.getGamesStatsMessage()
			}
		} else {
			c.Send <- &message{Action: "nogame"}
		}

	case "start":
		c.hub.Start(c)
	case "join":
		err := c.hub.join(m.Data, c)
		if err != nil {
			c.Send <- &message{Action: "join", Data: err.Error()}
		}
	case "begin":
		err := c.hub.begin(c)
		if err != nil {
			c.Send <- &message{Action: "error", Data: err.Error()}
		}
	case "play":
		i, err := strconv.Atoi(m.Data)
		if err != nil {
			c.Send <- &message{Action: "error", Data: "invalid card index send"}
		} else {
			err := c.hub.play(c, i)
			if err != nil {
				c.Send <- &message{Action: "error", Data: err.Error()}
			}
		}
	case "fetch":
		err := c.hub.fetchHand(c)
		if err != nil {
			c.Send <- &message{Action: "error", Data: err.Error()}
		}
	case "leave":
		c.hub.leave(c)
	case "restart":
		c.hub.restartGame(c)
	default:
		c.Send <- &message{Action: "error", Data: "invalid command"}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	newline = []byte{'\n'}
)

func (c *Client) waitForMsg() {
	defer func() {
		c.connection.Close()
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
		ticker.Stop()
		c.connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:

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

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				jsonMessage, err := json.Marshal(<-c.Send)

				if err != nil {
					return
				}
				w.Write(jsonMessage)
			}

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
