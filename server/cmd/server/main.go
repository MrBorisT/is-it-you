package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	addr      = ":8080"
	tickRate  = 20
	deltaTime = 1.0 / tickRate

	startX    = 80.0
	startY    = 220.0
	rowGap    = 50.0
	finishX   = 900.0
	walkSpeed = 80.0
	runSpeed  = 180.0
	hitRadius = 18.0
)

type ClientMessage struct {
	Type string `json:"type"`

	MoveRight bool `json:"move_right"`
	Running   bool `json:"running"`

	TargetX float64 `json:"target_x"`
	TargetY float64 `json:"target_y"`
}

type PlayerState struct {
	ID            string  `json:"id"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Running       bool    `json:"running"`
	ReachedFinish bool    `json:"reached_finish"`
	Alive         bool    `json:"alive"`
	HasBullet     bool    `json:"has_bullet"`
}

type ServerMessage struct {
	Type     string        `json:"type"`
	PlayerID string        `json:"player_id,omitempty"`
	Players  []PlayerState `json:"players,omitempty"`
	GameOver bool          `json:"game_over,omitempty"`
	WinnerID string        `json:"winner_id,omitempty"`
}

type Player struct {
	ID            string
	X             float64
	Y             float64
	MoveRight     bool
	Running       bool
	ReachedFinish bool
	Alive         bool
	HasBullet     bool
}

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
}

type Game struct {
	mu      sync.Mutex
	players map[string]*Player
	clients map[string]*Client

	gameOver bool
	winnerID string
}

var (
	nextID atomic.Int64

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	game = &Game{
		players: make(map[string]*Player),
		clients: make(map[string]*Client),
	}
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ws", wsHandler)

	go game.loop()

	log.Println("server listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	id := fmt.Sprintf("p%d", nextID.Add(1))

	client := &Client{
		id:   id,
		conn: conn,
		send: make(chan []byte, 16),
	}

	game.addClient(client)

	go client.writeLoop()
	client.readLoop()

	game.removeClient(id)
}

func (g *Game) addClient(c *Client) {
	g.mu.Lock()
	defer g.mu.Unlock()

	playerIndex := len(g.players)

	g.clients[c.id] = c
	g.players[c.id] = &Player{
		ID:        c.id,
		X:         startX,
		Y:         startY + float64(playerIndex)*rowGap,
		Alive:     true,
		HasBullet: true,
	}

	log.Println("client connected:", c.id)

	welcome := ServerMessage{
		Type:     "welcome",
		PlayerID: c.id,
	}

	data, err := json.Marshal(welcome)
	if err == nil {
		c.send <- data
	}
}

func (g *Game) removeClient(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if c, ok := g.clients[id]; ok {
		close(c.send)
		_ = c.conn.Close()
		delete(g.clients, id)
	}

	delete(g.players, id)

	log.Println("client disconnected:", id)
}

func (c *Client) readLoop() {
	defer c.conn.Close()

	for {
		var msg ClientMessage

		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Println("read error:", c.id, err)
			return
		}

		switch msg.Type {
		case "input":
			game.updateInput(c.id, msg)

		case "shoot":
			game.handleShoot(c.id, msg)

		default:
			log.Println("unknown message type:", msg.Type)
		}
	}
}

func (c *Client) writeLoop() {
	defer c.conn.Close()

	for data := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write error:", c.id, err)
			return
		}
	}
}

func (g *Game) updateInput(playerID string, msg ClientMessage) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.gameOver {
		return
	}

	player, ok := g.players[playerID]
	if !ok {
		return
	}

	if !player.Alive || player.ReachedFinish {
		player.MoveRight = false
		player.Running = false
		return
	}

	player.MoveRight = msg.MoveRight
	player.Running = msg.Running
}

func (g *Game) loop() {
	ticker := time.NewTicker(time.Second / tickRate)
	defer ticker.Stop()

	for range ticker.C {
		g.update()
		g.broadcastState()
	}
}

func (g *Game) update() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.gameOver {
		return
	}

	for _, player := range g.players {
		if !player.Alive || player.ReachedFinish {
			continue
		}

		if player.MoveRight {
			speed := walkSpeed
			if player.Running {
				speed = runSpeed
			}

			player.X += speed * deltaTime
		}

		if player.X >= finishX {
			player.X = finishX
			player.ReachedFinish = true
			g.gameOver = true
			g.winnerID = player.ID

			log.Println("winner by finish:", player.ID)
			return
		}
	}
}

func (g *Game) handleShoot(shooterID string, msg ClientMessage) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.gameOver {
		return
	}

	shooter, ok := g.players[shooterID]
	if !ok {
		return
	}

	if !shooter.Alive || shooter.ReachedFinish {
		return
	}

	if !shooter.HasBullet {
		return
	}

	shooter.HasBullet = false

	var hitPlayer *Player
	closestDistance := hitRadius * hitRadius

	for _, target := range g.players {
		if target.ID == shooterID {
			continue
		}

		if !target.Alive || target.ReachedFinish {
			continue
		}

		distance := distanceSquared(target.X, target.Y, msg.TargetX, msg.TargetY)
		if distance <= closestDistance {
			closestDistance = distance
			hitPlayer = target
		}
	}

	if hitPlayer == nil {
		log.Println("shot missed by:", shooterID)
		return
	}

	hitPlayer.Alive = false
	hitPlayer.MoveRight = false
	hitPlayer.Running = false

	log.Println("player", shooterID, "shot", hitPlayer.ID)

	g.checkWinAfterKill(shooterID)
}

func (g *Game) checkWinAfterKill(killerID string) {
	aliveCount := 0
	lastAliveID := ""

	for _, player := range g.players {
		if player.Alive && !player.ReachedFinish {
			aliveCount++
			lastAliveID = player.ID
		}
	}

	if aliveCount == 1 {
		g.gameOver = true
		g.winnerID = lastAliveID
		log.Println("winner by elimination:", lastAliveID)
		return
	}

	if aliveCount == 0 {
		g.gameOver = true
		g.winnerID = killerID
		log.Println("winner fallback:", killerID)
	}
}

func distanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

func (g *Game) broadcastState() {
	g.mu.Lock()
	defer g.mu.Unlock()

	players := make([]PlayerState, 0, len(g.players))

	for _, player := range g.players {
		players = append(players, PlayerState{
			ID:            player.ID,
			X:             player.X,
			Y:             player.Y,
			Running:       player.Running,
			ReachedFinish: player.ReachedFinish,
			Alive:         player.Alive,
			HasBullet:     player.HasBullet,
		})
	}

	msg := ServerMessage{
		Type:     "state",
		Players:  players,
		GameOver: g.gameOver,
		WinnerID: g.winnerID,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal state error:", err)
		return
	}

	for id, client := range g.clients {
		select {
		case client.send <- data:
		default:
			log.Println("client send buffer full, dropping:", id)
		}
	}
}
