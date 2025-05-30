package logic

import (
	"connect4/server/internal/game/model"
)

const (
	target = 4
)

type Board [model.Column][model.Row]uint8

type Location struct {
	Column int
	Row    int
}

type FourInARow [target]Location

func diagonal4InRow(b Board) *FourInARow {
	l := downRightDiagonal4InRow(b)
	if l != nil {
		return l
	}

	return upRightDiagonal4InRow(b)
}

func upRightDiagonal4InRow(b Board) *FourInARow {
	f := new(FourInARow)
	for c := 0; c < model.Column; c++ {
		if c+target-1 >= model.Column {
			break
		}

		for r := 0; r < model.Row; r++ {
			if r+target-1 >= model.Row {
				break
			}

			loc := Location{
				Column: c,
				Row:    r,
			}

			p := b[c][r]
			if p == model.NoPlayer {
				continue
			}

			len := 0
			for loc.Column < model.Column && loc.Row < model.Row && p == b[loc.Column][loc.Row] {
				// Create new location
				nLoc := Location{
					Column: loc.Column,
					Row:    loc.Row,
				}
				// Set location in FourInARow obj
				f[len] = nLoc
				// Increment length
				len++
				// Set current player at location
				p = b[nLoc.Column][nLoc.Row]
				// Traverse right diagonally
				loc.Column++
				loc.Row++
				if len == target {
					return f
				}
			}
		}
	}

	return nil
}

func downRightDiagonal4InRow(b Board) *FourInARow {
	f := new(FourInARow)
	for c := 0; c+target-1 < model.Column; c++ {
		for r := model.Row - 1; r-target+1 >= 0; r-- {
			loc := Location{
				Column: c,
				Row:    r,
			}

			p := b[c][r]
			if p == model.NoPlayer {
				continue
			}

			len := 0
			for loc.Column < model.Column && loc.Row < model.Row && p == b[loc.Column][loc.Row] {
				// Create new location
				nLoc := Location{
					Column: loc.Column,
					Row:    loc.Row,
				}
				// Set location in FourInARow obj
				f[len] = nLoc
				// Increment length
				len++
				// Set current player at location
				p = b[nLoc.Column][nLoc.Row]
				// Traverse left diagonally
				loc.Column++
				loc.Row--
				if len == target {
					return f
				}
			}
		}
	}

	return nil
}

func vertical4InRow(b Board) *FourInARow {
	f := new(FourInARow)
	for c := range model.Column {
		for p1 := 0; p1+target-1 < model.Row; {
			start := b[c][p1]
			if start == model.NoPlayer {
				p1++
				continue
			}

			for p2 := p1; p2 < p1+target; p2++ {
				if b[c][p2] != start {
					p1 = p2
					break
				}

				i := p2 - p1
				f[i] = Location{
					Column: c,
					Row:    p2,
				}

				if p2 == p1+target-1 {
					return f
				}
			}
		}
	}

	return nil
}

func horizontal4InRow(b Board) *FourInARow {
	f := new(FourInARow)
	for r := range model.Row {
		for p1 := 0; p1+target-1 < model.Column; {
			start := b[p1][r]
			if start == model.NoPlayer {
				p1++
				continue
			}

			for p2 := p1; p2 < p1+target; p2++ {
				if b[p2][r] != start {
					p1 = p2
					break
				}

				i := p2 - p1
				f[i] = Location{
					Column: p2,
					Row:    r,
				}

				if p2 == p1+target-1 {
					return f
				}
			}
		}
	}

	return nil
}

func Has4InARow(b Board) *FourInARow {
	diag := diagonal4InRow(b)
	if diag != nil {
		return diag
	}

	vert := vertical4InRow(b)
	if vert != nil {
		return vert
	}

	hrz := horizontal4InRow(b)
	if hrz != nil {
		return hrz
	}

	return nil
}

func (l1 Location) Compare(l2 Location) int {
	if l1.Column == l2.Column && l1.Row == l2.Row {
		return 0
	}

	return -1
}

type direction struct {
	dx int // change in col
	dy int // change in row
}

func upRightDirection() direction {
	return direction{
		dx: 1,
		dy: 1,
	}
}

func downRightDirection() direction {
	return direction{
		dx: 1,
		dy: -1,
	}
}
