package logic

import (
	"connect4/server/internal/game/model"
)

type board [model.Column][model.Row]uint8

type Location struct {
	Column int
	Row    int
}

type FourInARow [4]Location

func diagonal4InRow(b board) *FourInARow {
	return nil
}

func vertical4InRow(b board) *FourInARow {
	f := new(FourInARow)
	target := 4
	for c := 0; c < model.Column; c++ {
		for p1 := 0; p1+3 < model.Row; {
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

func horizontal4InRow(b board) *FourInARow {
	return nil
}

func Has4InARow(g model.Game) *FourInARow {
	b := g.GetBoard()
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
