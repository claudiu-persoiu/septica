package game

import (
	"github.com/gorilla/websocket"
	"time"
	"log"
	"bytes"
	"net/http"
	"encoding/json"
	"fmt"
)

type Client struct {
	connection *websocket.Conn
	Send       chan *Message
	hub        *Hub
	gameKey    string
}

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data,omitempty"`
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request, hub *Hub) *Client {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	client := &Client{connection: conn, Send: make(chan *Message, 256), hub: hub}

	go client.writePump()
	go client.readPump()

	return client
}

func (c *Client) processMessage(message Message) {

	switch message.Action {
	case "start":
		c.hub.Start <- c
	case "join":
		err := c.hub.Join(message.Data, c)
		if err != nil {
			c.Send <- &Message{Action: "join", Data: err.Error()}
		} else {
			c.Send <- &Message{Action: "join", Data: "wait"}
		}
	case "begin":
		c.hub.Begin(c)
	default:
		c.Send <- &Message{Action: "error", Data: "invalid command"}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (c *Client) readPump() {
	defer func() {
		//c.hub.unregister <- c
		c.connection.Close()
	}()
	c.connection.SetReadLimit(maxMessageSize)
	c.connection.SetReadDeadline(time.Now().Add(pongWait))
	c.connection.SetPongHandler(func(string) error { c.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		fmt.Println(string(message))

		var obj Message
		if err := json.Unmarshal(message, &obj); err == nil {
			c.processMessage(obj)
		}

	}
}

func (c *Client) writePump() {
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
