package gameflow

import (
	"errors"
	"log"
	"sync"
)

const (
	GameNotReady     = "game not ready to play"
	PlayerNotInAGame = "player is not in a game"
	FailedAction     = "failed to do action"
)

type OpponentPlayed struct {
	Column int
}

var observerMutex sync.Mutex

func (p *Player) PlayPiece(col int) error {
	if p.game == nil {
		return errors.New(PlayerNotInAGame)
	} else if !p.game.readyToPlay {
		return errors.New(GameNotReady)
	}

	if !p.game.handleAction(p, doPlay{Column: col}) {
		return errors.New(FailedAction)
	}

	return nil
}

func (p *Player) GetOpponentUsername() *string {
	if p.game == nil {
		return nil
	}

	opp := p.game.getOpponent(p)
	if opp == nil {
		return nil
	}

	return &opp.Username
}

func (p *Player) GetGameWinner() string {
	if p.game == nil {
		return ""
	}

	return p.game.getWinner()
}

func (p *Player) GetStartingPlayer() string {
	if p.game == nil {
		return ""
	}

	return p.game.getStartingPlayer()
}

func (p *Player) IsYourTurn() bool {
	if p.game == nil {
		return false
	}

	return p.game.isPlayersTurn(p)
}

func (p *Player) GetPlays() []int {
	if p.game == nil {
		return []int{}
	}

	return p.game.getPlays()
}

func (p *Player) GetBoard() uint64 {
	if p.game == nil {
		return 0
	}

	return p.game.getBoard()
}

func (p *Player) FindNewGame() {
	if p.game == nil {
		return
	}

	p.game.leaveGame(p)

	JoinGame(p)
}

func (p *Player) handleAction(action interface{}) {
	observerMutex.Lock()
	defer observerMutex.Unlock()

	if p.observer == nil {
		log.Printf("%s observer is nil", p.Username)
		return
	}

	f := *p.observer
	f(action)
}
