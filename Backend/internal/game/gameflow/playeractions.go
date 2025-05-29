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
