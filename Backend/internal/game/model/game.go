package model

import "errors"

const (
	Column = 7
	Row    = 6
)

type Game struct {
	board [Column][Row]int
}

// Drop a piece in the desired column.
// Takes a player p (false for p1, true for p2) and
// the desired column c to drop a piece in.
func (g Game) dropPiece(p bool, c int) error {
	if c < 0 || c > Column-1 {
		return errors.New("Column is out of range")
	}

	spot := -1
	for i := 0; i < len(g.board[c]); i++ {
		if g.board[c][i] == 0 {
			spot = i
			break
		}
	}

	if spot == -1 {
		return errors.New("Column is full")
	}

	if p {
		g.board[c][spot] = 1
	} else {
		g.board[c][spot] = 2
	}

	return nil
}
