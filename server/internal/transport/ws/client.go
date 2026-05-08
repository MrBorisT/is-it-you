package ws

import (
	"encoding/json"
	"log"

	"github.com/MrBorisT/is-it-you-server/internal/protocol"
	"github.com/gorilla/websocket"
)

type Client struct {
	id      string
	conn    *websocket.Conn
	send    chan []byte
	handler *Handler
}

func NewClient(id string, conn *websocket.Conn, handler *Handler) *Client {
	return &Client{
		id:      id,
		conn:    conn,
		send:    make(chan []byte, 16),
		handler: handler,
	}
}

func (c *Client) ReadLoop() {
	defer c.handler.RemoveClient(c.id)

	for {
		var msg protocol.ClientMessage

		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Println("read error:", c.id, err)
			return
		}

		switch msg.Type {
		case "input":
			c.handler.game.UpdateInput(c.id, msg)

		case "shoot":
			c.handler.game.HandleShoot(c.id, msg)

		case "restart":
			c.handler.game.RestartRound()

		default:
			log.Println("unknown message type:", msg.Type)
		}
	}
}

func (c *Client) WriteLoop() {
	defer c.conn.Close()

	for data := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write error:", c.id, err)
			return
		}
	}
}

func (c *Client) SendMessage(msg protocol.ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal message error:", err)
		return
	}

	select {
	case c.send <- data:
	default:
		log.Println("client send buffer full, dropping:", c.id)
	}
}

func (c *Client) Close() {
	close(c.send)
	_ = c.conn.Close()
}
