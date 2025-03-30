package model

import (
	"testing"
)

func TestGame_DropPieceNegativeColumn(t *testing.T) {
	g := new(Game)

	e := g.DropPiece(false, -1)
	if e == nil {
		t.Error("Expected error to throw for column out of bounds")
	}
}

func TestGame_DropPieceLargeColumn(t *testing.T) {
	g := new(Game)

	e := g.DropPiece(false, Column+1239843)
	if e == nil {
		t.Error("Expected error to throw for column out of bounds")
	}
}

func (g *Game) fillColumn(c int) {
	for i := range Row {
		g.board[c][i] = 1
	}
}

func TestGame_DropPieceColumnFull(t *testing.T) {
	g := new(Game)
	for i := range Column {
		g.fillColumn(i)

		e := g.DropPiece(false, i)
		if e == nil {
			t.Error("Expected error to throw when column is full")
		}
	}
}

func TestGame_DropPieceValid(t *testing.T) {
	g := new(Game)
	for i := range Column {
		for j := range Row {
			e := g.DropPiece((i+j)%2 == 0, i)
			if e != nil {
				t.Error("Dropping piece in valid state should not throw error")
			}
		}

		e := g.DropPiece(false, i)
		if e == nil {
			t.Error("Column should be filled")
		}
	}
}
