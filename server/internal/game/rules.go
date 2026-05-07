package game

func (g *Game) updateLocked() {
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

	var hitPlayer *Player
	closestDistance := HitRadius * HitRadius

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
		}
	}

	if hitPlayer == nil {
		return
	}

	hitPlayer.Alive = false
	hitPlayer.MoveRight = false
	hitPlayer.Running = false

	g.checkWinAfterKillLocked(shooterID)
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
