package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/MrBorisT/is-it-you-server/internal/protocol"
)

const (
	TickRate  = 20
	DeltaTime = 1.0 / TickRate

	StartX    = 80.0
	StartY    = 220.0
	RowGap    = 50.0
	FinishX   = 900.0
	WalkSpeed = 80.0
	RunSpeed  = 180.0
	HitRadius = 18.0

	NPCCount     = 25
	NPCWalkSpeed = 55.0
	NPCMinStep   = 15.0
	NPCMaxStep   = 70.0
	NPCMinWait   = 0.3
	NPCMaxWait   = 1.8
	NPCMinY      = 140.0
	NPCMaxY      = 420.0
)

type Game struct {
	mu sync.Mutex

	players map[string]*Player
	npcs    map[string]*NPC

	gameOver bool
	winnerID string

	rng *rand.Rand
}

func NewGame() *Game {
	g := &Game{
		players: make(map[string]*Player),
		npcs:    make(map[string]*NPC),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	g.spawnNPCs()

	return g
}

func (g *Game) spawnNPCs() {
	for i := 0; i < NPCCount; i++ {
		id := fmt.Sprintf("n%d", i+1)

		y := NPCMinY + g.rng.Float64()*(NPCMaxY-NPCMinY)

		g.npcs[id] = &NPC{
			ID:        id,
			X:         StartX,
			Y:         y,
			Alive:     true,
			TargetX:   StartX,
			WaitTimer: randomRange(g.rng, NPCMinWait, NPCMaxWait),
		}
	}
}

func (g *Game) AddPlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	playerIndex := len(g.players)

	g.players[id] = &Player{
		ID:        id,
		X:         StartX,
		Y:         StartY + float64(playerIndex)*RowGap,
		Alive:     true,
		HasBullet: true,
	}
}

func (g *Game) RemovePlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.players, id)
}

func (g *Game) UpdateInput(playerID string, msg protocol.ClientMessage) {
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
