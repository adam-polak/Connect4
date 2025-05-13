package logic

import "connect4/server/internal/game/model"

type board [model.Column][model.Row]uint8

type location struct {
	column int
	row    int
}

func diagonal4InRow(b board) *location {
	return nil
}

func vertical4InRow(b board) *location {
	return nil
}

func horizontal4InRow(b board) *location {
	return nil
}

func checkForWin(g model.Game) bool {
	b := g.GetBoard()
	diag := diagonal4InRow(b)
	if diag != nil {
		return true
	}

	vert := vertical4InRow(b)
	if vert != nil {
		return true
	}

	hrz := horizontal4InRow(b)
	if hrz != nil {
		return true
	}

	return false
}
