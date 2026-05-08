package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/MrBorisT/is-it-you/server/internal/config"
	"github.com/MrBorisT/is-it-you/server/internal/protocol"
)

const (
	TickRate = 20
)

type Game struct {
	mu sync.Mutex

	cfg config.GameConfig

	players map[string]*Player
	npcs    map[string]*NPC

	gameOver bool
	winnerID string

	rng *rand.Rand
}

func NewGame(cfg config.GameConfig) *Game {
	g := &Game{
		cfg:     cfg,
		players: make(map[string]*Player),
		npcs:    make(map[string]*NPC),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	g.spawnNPCs()

	return g
}

func (g *Game) spawnNPCs() {
	for i := 0; i < g.cfg.NPCCount; i++ {
		id := fmt.Sprintf("n%d", i+1)

		y := g.cfg.NPCMinY + g.rng.Float64()*(g.cfg.NPCMaxY-g.cfg.NPCMinY)

		g.npcs[id] = &NPC{
			ID:        id,
			X:         g.cfg.StartX,
			Y:         y,
			Alive:     true,
			TargetX:   g.cfg.StartX,
			WaitTimer: randomRange(g.rng, g.cfg.NPCMinWait, g.cfg.NPCMaxWait),
		}
	}
}

func (g *Game) AddPlayer(id string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.players) >= g.cfg.MaxPlayers {
		return false
	}

	playerIndex := len(g.players)

	g.players[id] = &Player{
		ID:        id,
		X:         g.cfg.StartX,
		Y:         g.cfg.StartY + float64(playerIndex)*g.cfg.RowGap,
		AimX:      g.cfg.StartX,
		AimY:      g.cfg.StartY + float64(playerIndex)*g.cfg.RowGap,
		Alive:     true,
		HasBullet: true,
	}

	return true
}

func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.players, id)
}

func (g *Game) UpdateInput(playerID string, msg protocol.ClientMessage) {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, ok := g.players[playerID]
	if !ok {
		return
	}

	if g.gameOver {
		return
	}

	player.AimX = msg.AimX
	player.AimY = msg.AimY

	if !player.Alive || player.ReachedFinish {
		player.MoveRight = false
		player.Running = false
		return
	}

	player.MoveRight = msg.MoveRight
	player.Running = msg.Running
}

func (g *Game) HandleShoot(shooterID string, msg protocol.ClientMessage) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.handleShootLocked(shooterID, msg.TargetX, msg.TargetY)
}

func (g *Game) Loop(onState func(protocol.ServerMessage)) {
	ticker := time.NewTicker(time.Second / TickRate)
	defer ticker.Stop()

	for range ticker.C {
		state := g.Tick()
		onState(state)
	}
}

func (g *Game) Tick() protocol.ServerMessage {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.gameOver {
		g.updateLocked()
	}

	return g.stateLocked()
}

func (g *Game) stateLocked() protocol.ServerMessage {
	players := make([]protocol.PlayerState, 0, len(g.players))

	for _, player := range g.players {
		players = append(players, protocol.PlayerState{
			ID:            player.ID,
			X:             player.X,
			Y:             player.Y,
			Running:       player.Running,
			ReachedFinish: player.ReachedFinish,
			Alive:         player.Alive,
			HasBullet:     player.HasBullet,
			AimX:          player.AimX,
			AimY:          player.AimY,
		})
	}

	npcs := make([]protocol.NPCState, 0, len(g.npcs))

	for _, npc := range g.npcs {
		npcs = append(npcs, protocol.NPCState{
			ID:            npc.ID,
			X:             npc.X,
			Y:             npc.Y,
			Alive:         npc.Alive,
			ReachedFinish: npc.ReachedFinish,
		})
	}

	return protocol.ServerMessage{
		Type:     "state",
		Players:  players,
		NPCs:     npcs,
		GameOver: g.gameOver,
		WinnerID: g.winnerID,
	}
}

func (g *Game) PlayerCount() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	return len(g.players)
}

func (g *Game) RestartRound() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.gameOver {
		return false
	}

	g.gameOver = false
	g.winnerID = ""

	i := 0
	for _, player := range g.players {
		player.X = g.cfg.StartX
		player.Y = g.cfg.StartY + float64(i)*g.cfg.RowGap

		player.MoveRight = false
		player.Running = false

		player.AimX = player.X
		player.AimY = player.Y

		player.ReachedFinish = false
		player.Alive = true
		player.HasBullet = true

		i++
	}

	// reset NPCs
	g.npcs = make(map[string]*NPC)
	g.spawnNPCs()

	return true
}

func (g *Game) deltaTime() float64 {
	return 1.0 / float64(g.cfg.TickRate)
}
