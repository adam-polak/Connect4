package logic

import (
	"connect4/server/internal/game/model"
	"errors"
	"math/rand"
)

type Connect4Engine struct {
	winner   *string
	round    uint
	game     *model.Game
	p1       string
	p2       string
	cur_turn bool
}

func NewConnect4Engine(p1 string, p2 string) Connect4Engine {
	g := new(model.Game)
	cur := rand.Int()%2 == 0
	return Connect4Engine{
		winner:   nil,
		round:    1,
		game:     g,
		p1:       p1,
		p2:       p2,
		cur_turn: cur,
	}
}

func (c *Connect4Engine) GetWinner() *string {
	return c.winner
}

func (c *Connect4Engine) CurrentRound() uint {
	return c.round
}

func (c *Connect4Engine) Board() uint64 {
	return c.game.Encode()
}

func (c *Connect4Engine) DropPiece(p *string, col int) error {
	if c.winner != nil {
		return errors.New("game is already over")
	} else if c.cur_turn && *p != c.p1 {
		return errors.New("not player two's turn")
	} else if !c.cur_turn && *p != c.p2 {
		return errors.New("not player one's turn")
	}

	// try to drop piece
	err := c.game.DropPiece(c.cur_turn, col)
	if err != nil {
		return err
	}

	// check for four ina row
	f := Has4InARow(*c.game)
	if f != nil {
		// declare winner
		c.winner = p
	}

	// advance to next turn
	c.cur_turn = !c.cur_turn
	c.round++

	return nil
}
