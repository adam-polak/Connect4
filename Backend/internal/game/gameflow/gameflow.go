package gameflow

import (
	"connect4/server/internal/game/logic"
	"sync"
)

type GameOrchestrator struct {
	player1     *Player
	player2     *Player
	readyToPlay bool
	engine      logic.Connect4Engine
}

type gameNode struct {
	game *GameOrchestrator
	next *gameNode
	prev *gameNode
}

var availableGames = &gameNode{
	game: nil,
	next: nil,
	prev: nil,
}

type doPlay struct {
	Column int
}

type GameOver struct {
	Winner string
}

type GameReady struct{}

var availableGamesMutex sync.Mutex

func init() {
	availableGames.next = availableGames
	availableGames.prev = availableGames
}

func JoinGame(p *Player) {
	availableGamesMutex.Lock()
	defer availableGamesMutex.Unlock()

	var g *GameOrchestrator
	if availableGames.next == availableGames {
		// create new game
		g = &GameOrchestrator{
			player1: p,
		}

		// add game to available games list
		availableGames.next = &gameNode{
			game: g,
			next: availableGames,
			prev: availableGames,
		}

		availableGames.prev = availableGames.next
	} else {
		// get available game node
		n := availableGames.prev
		// remove node from available list
		availableGames.prev = n.prev
		availableGames.prev.next = availableGames
		// set game to the available game
		g = n.game
		// set player 2
		g.player2 = p
		// set game to ready state
		g.readyToPlay = true
		// notify players
		g.notifyPlayers(GameReady{})
		// create game engine
		g.engine = logic.NewConnect4Engine(g.player1.Username, g.player2.Username)
	}

	p.game = g
}

func (g *GameOrchestrator) getOpponent(p *Player) *Player {
	if p == g.player1 {
		return g.player2
	}

	return g.player1
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
