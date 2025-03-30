package model

import "errors"

const (
	Column     = 7
	Row        = 6
	RangeError = "Column is out of range"
	FullError  = "Column is full"
)

type Game struct {
	board [Column][Row]uint8
}

// Drop a piece in the desired column.
// Takes a player p (false for p1, true for p2) and
// the desired column c to drop a piece in.
func (g *Game) DropPiece(p bool, c int) error {
	full, e := g.ColumnFull(c)
	if e != nil {
		return e
	} else if full {
		return errors.New(FullError)
	}

	spot := -1
	for i := 0; i < len(g.board[c]); i++ {
		if g.board[c][i] == 0 {
			spot = i
			break
		}
	}

	if p {
		g.board[c][spot] = 1
	} else {
		g.board[c][spot] = 2
	}

	return nil
}

// Return true if the column of the game board is full.
func (g Game) ColumnFull(c int) (bool, error) {
	if c < 0 || c > Column-1 {
		return false, errors.New(RangeError)
	}

	return g.board[c][Row-1] != 0, nil
}
