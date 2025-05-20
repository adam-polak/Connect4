package gameflow

import "connect4/server/internal/game/player"

type GameOrchestrator struct {
	p1 player.Player
	p2 player.Player
}
