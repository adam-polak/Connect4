package gameflow

import "errors"

const (
	GameNotReady     = "game not ready to play"
	PlayerNotInAGame = "player is not in a game"
	FailedAction     = "failed to do action"
)

type PlayerObserver interface {
	HandleAction(interface{})
}

type OpponentPlayed struct {
	Column int
}

func (p *Player) PlayPiece(col int) error {
	if p.game == nil || !p.game.readyToPlay {
		return errors.New(PlayerNotInAGame)
	} else if !p.game.readyToPlay {
		return errors.New(GameNotReady)
	}

	if !p.game.handleAction(p, doPlay{
		Column: col,
	}) {
		return errors.New(FailedAction)
	}

	return nil
}

func (p *Player) notifyObservers(action interface{}) {
	for i := range p.observers {
		p.observers[i].HandleAction(action)
	}
}

func (p *Player) RegisterObserver(o PlayerObserver) {
	p.observers = append(p.observers, o)
}

func (p *Player) handleAction(action interface{}) {
	p.notifyObservers(action)
}
