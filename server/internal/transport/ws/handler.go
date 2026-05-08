package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/MrBorisT/is-it-you/server/internal/game"
	"github.com/MrBorisT/is-it-you/server/internal/protocol"
	"github.com/gorilla/websocket"
)

type Handler struct {
	game *game.Game

	mu      sync.Mutex
	clients map[string]*Client
	nextID  atomic.Int64

	upgrader websocket.Upgrader
}

func NewHandler(g *game.Game) *Handler {
	return &Handler{
		game:    g,
		clients: make(map[string]*Client),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	id := fmt.Sprintf("p%d", h.nextID.Add(1))
	client := NewClient(id, conn, h)

	if ok := h.AddClient(client); !ok {
		return
	}

	go client.WriteLoop()
	client.ReadLoop()
}

func (h *Handler) AddClient(c *Client) bool {
	if ok := h.game.AddPlayer(c.id); !ok {
		c.SendMessage(protocol.ServerMessage{
			Type:    "error",
			Message: "room full",
		})
		c.Close()
		return false
	}

	h.mu.Lock()
	h.clients[c.id] = c
	h.mu.Unlock()

	log.Println("client connected:", c.id)

	c.SendMessage(protocol.ServerMessage{
		Type:     "welcome",
		PlayerID: c.id,
	})

	return true
}

func (h *Handler) RemoveClient(id string) {
	h.mu.Lock()

	client, ok := h.clients[id]
	if ok {
		delete(h.clients, id)
	}

	h.mu.Unlock()

	if ok {
		client.Close()
	}

	h.game.RemovePlayer(id)

	log.Println("client disconnected:", id)
}

func (h *Handler) Broadcast(msg protocol.ServerMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.clients {
		client.SendMessage(msg)
	}
}
