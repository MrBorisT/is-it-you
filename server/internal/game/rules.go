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
			speed := g.cfg.WalkSpeed
			if player.Running {
				speed = g.cfg.RunSpeed
			}

			player.X += speed * g.deltaTime()
		}

		if player.X >= g.cfg.FinishX {
			player.X = g.cfg.FinishX
			player.ReachedFinish = true
			g.phase = PhaseFinished
			g.winnerID = player.ID

			for _, p := range g.players {
				p.MoveRight = false
				p.Running = false
			}

			return
		}
	}
}

func (g *Game) updateNPCsLocked() {
	for _, npc := range g.npcs {
		if !npc.Alive || npc.ReachedFinish {
			continue
		}

		if npc.X >= g.cfg.FinishX {
			npc.X = g.cfg.FinishX
			npc.ReachedFinish = true
			continue
		}

		if !npc.Moving {
			npc.WaitTimer -= g.deltaTime()

			if npc.WaitTimer <= 0 {
				step := randomRange(g.rng, g.cfg.NPCMinStep, g.cfg.NPCMaxStep)
				npc.TargetX = npc.X + step
				if npc.TargetX > g.cfg.FinishX {
					npc.TargetX = g.cfg.FinishX
				}
				npc.Moving = true
			}

			continue
		}

		npc.X += g.cfg.NPCWalkSpeed * g.deltaTime()

		if npc.X >= npc.TargetX {
			npc.X = npc.TargetX
			npc.Moving = false
			npc.WaitTimer = randomRange(g.rng, g.cfg.NPCMinWait, g.cfg.NPCMaxWait)
		}

		if npc.X >= g.cfg.FinishX {
			npc.X = g.cfg.FinishX
			npc.ReachedFinish = true
			npc.Moving = false
		}
	}
}

func (g *Game) handleShootLocked(shooterID string, targetX, targetY float64) {
	if g.phase != PhaseRunning {
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

	closestDistance := g.cfg.HitRadius * g.cfg.HitRadius

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

func distanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

func randomRange(rng *rand.Rand, min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}
