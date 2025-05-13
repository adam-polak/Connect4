package logic

import "connect4/server/internal/game/model"

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
