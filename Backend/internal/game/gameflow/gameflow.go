package gameflow

import "connect4/server/internal/game/logic"

type GameOrchestrator struct {
	player1     *Player
	player2     *Player
	readyToPlay bool
	engine      logic.Connect4Engine
}

type gameNode struct {
	game *GameOrchestrator
	next *gameNode
}

var availableGames = &gameNode{
	game: nil,
	next: nil,
}

func init() {
	availableGames.next = availableGames
}

func GetGameOrchestrator(p *Player) *GameOrchestrator {
	var g *GameOrchestrator
	if availableGames.next == availableGames {
		// create new game
		g = &GameOrchestrator{
			player1: p,
		}

		// add game to available games list
		availableGames.next = &gameNode{
			game: g,
			next: availableGames.next,
		}
	} else {
		// get available game node
		n := availableGames.next
		// remove node from available list
		availableGames.next = n.next
		// set game to the available game
		g = n.game
		// set player 2
		g.player2 = p
		// set game to ready state
		g.readyToPlay = true
		// create game engine
		g.engine = logic.NewConnect4Engine(g.player1.Username, g.player2.Username)
	}

	p.game = g

	return g
}

func (g *GameOrchestrator) getOpponent(p *Player) *Player {
	if p == g.player1 {
		return g.player2
	}

	return g.player1
}

type doPlay struct {
	Column int
}

type GameOver struct {
	Winner string
}

func (g *GameOrchestrator) notifyPlayers(action interface{}) {
	g.player1.handleAction(action)
	g.player2.handleAction(action)
}

func (g *GameOrchestrator) handleAction(p *Player, action interface{}) bool {
	if g.player1 != p && g.player2 != p {
		panic("player is not in game")
	}

	switch v := action.(type) {
	case doPlay:
		if g.engine.GetWinner() != nil {
			return false
		}

		err := g.engine.DropPiece(&p.Username, v.Column)
		if err != nil {
			return false
		} else {
			g.getOpponent(p).handleAction(OpponentPlayed(v))
			if g.engine.GetWinner() != nil {
				g.notifyPlayers(GameOver{Winner: *g.engine.GetWinner()})
			}

			return true
		}
	default:
		panic("action not supported")
	}
}
