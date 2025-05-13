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
	for c := range model.Column {

		// * note use sliding window instead
		for r := 0; r < model.Row-4; r++ {
			start := b[c][r]
			if start == model.NoPlayer {
				break
			}

			len := 0
			for r2 := r; r2 < r+4; r2++ {
				if b[c][r2] != start || len == 4 {
					break
				}

				f[r2-r] = Location{
					Column: c,
					Row:    r2,
				}

				len++
			}

			if len == 4 {
				return f
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
