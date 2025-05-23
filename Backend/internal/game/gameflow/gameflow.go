package gameflow

type GameOrchestrator struct {
	player1 *Player
	player2 *Player
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
		// set game to the available game
		g = n.game
		// set player 2
		g.player2 = p
		// remove node from available list
		availableGames.next = n.next
	}

	p.game = g

	return g
}
