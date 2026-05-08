package game

import "math/rand"

func (g *Game) updateLocked() {
	g.updatePlayersLocked()
	g.updateNPCsLocked()
}

func (g *Game) updatePlayersLocked() {
	for _, player := range g.players {
		if !player.Alive || player.ReachedFinish {
			continue
		}

		if player.MoveRight {
			speed := WalkSpeed
			if player.Running {
				speed = RunSpeed
			}

			player.X += speed * DeltaTime
		}

		if player.X >= FinishX {
			player.X = FinishX
			player.ReachedFinish = true
			g.gameOver = true
			g.winnerID = player.ID
			return
		}
	}
}

func (g *Game) updateNPCsLocked() {
	for _, npc := range g.npcs {
		if !npc.Alive || npc.ReachedFinish {
			continue
		}

		if npc.X >= FinishX {
			npc.X = FinishX
			npc.ReachedFinish = true
			continue
		}

		if !npc.Moving {
			npc.WaitTimer -= DeltaTime

			if npc.WaitTimer <= 0 {
				step := randomRange(g.rng, NPCMinStep, NPCMaxStep)
				npc.TargetX = npc.X + step
				if npc.TargetX > FinishX {
					npc.TargetX = FinishX
				}
				npc.Moving = true
			}

			continue
		}

		npc.X += NPCWalkSpeed * DeltaTime

		if npc.X >= npc.TargetX {
			npc.X = npc.TargetX
			npc.Moving = false
			npc.WaitTimer = randomRange(g.rng, NPCMinWait, NPCMaxWait)
		}

		if npc.X >= FinishX {
			npc.X = FinishX
			npc.ReachedFinish = true
			npc.Moving = false
		}
	}
}

func (g *Game) handleShootLocked(shooterID string, targetX, targetY float64) {
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

	closestDistance := HitRadius * HitRadius

	var hitPlayer *Player
	var hitNPC *NPC

	for _, target := range g.players {
		if target.ID == shooterID {
			continue
		}

		if !target.Alive || target.ReachedFinish {
			continue
		}

		distance := distanceSquared(target.X, target.Y, targetX, targetY)
		if distance <= closestDistance {
			closestDistance = distance
			hitPlayer = target
			hitNPC = nil
		}
	}

	for _, npc := range g.npcs {
		if !npc.Alive || npc.ReachedFinish {
			continue
		}

		distance := distanceSquared(npc.X, npc.Y, targetX, targetY)
		if distance <= closestDistance {
			closestDistance = distance
			hitNPC = npc
			hitPlayer = nil
		}
	}

	if hitPlayer != nil {
		hitPlayer.Alive = false
		hitPlayer.MoveRight = false
		hitPlayer.Running = false
		return
	}

	if hitNPC != nil {
		hitNPC.Alive = false
		hitNPC.Moving = false
		return
	}
}

func (g *Game) checkWinAfterKillLocked(killerID string) {
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
		return
	}

	if aliveCount == 0 {
		g.gameOver = true
		g.winnerID = killerID
	}
}

func distanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

func randomRange(rng *rand.Rand, min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}
