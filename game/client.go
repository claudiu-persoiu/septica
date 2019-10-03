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
	Send       chan *message
	hub        *Hub
	cards      []*card
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
	case "start":
		c.hub.Start(c)
	case "join":
		err := c.hub.join(m.Data, c)
		if err != nil {
			c.Send <- &message{Action: "join", Data: err.Error()}
		} else {
			c.Send <- &message{Action: "join", Data: "wait"}
		}
	case "begin":
		c.hub.begin(c)
	case "play":
		i, err := strconv.Atoi(m.Data)
		if err != nil {
			c.Send <- &message{Action: "error", Data: "invalid card index send"}
		} else {
			c.hub.play(c, i)
		}
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
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (c *Client) waitForMsg() {
	defer c.connection.Close()
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
			log.Println("Error parcing message:")
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
