package logic

import (
	"connect4/server/internal/game/model"
	"testing"
)

func Test_BasicVertical(t *testing.T) {
	for c := 0; c < model.Column; c++ {
		g := new(model.Game)
		for range 4 {
			g.DropPiece(c%2 == 0, c)
		}

		f := vertical4InRow(g.GetBoard())
		if f == nil {
			t.Error("Should detect vertical four in a row")
			return
		}

		for r := range 4 {
			if f[r].Column != c {
				t.Errorf("Expected column %d was %d", c, f[r].Column)
				return
			}

			if f[r].Row != r {
				t.Errorf("Expected row %d was %d", r, f[r].Row)
				return
			}
		}
	}
}

func Test_Last4Vertical(t *testing.T) {
	for c := range model.Column {
		g := new(model.Game)
		for r := range model.Row {
			if r < 2 {
				g.DropPiece(c%2 != 0, c)
			} else {
				g.DropPiece(c%2 == 0, c)
			}
		}

		f := vertical4InRow(g.GetBoard())
		if f == nil {
			t.Error("Should detect vertical four in a row")
			return
		}

		for r := range 4 {
			if f[r].Column != c {
				t.Errorf("Expected column %d was %d", c, f[r].Column)
				return
			}

			if f[r].Row != r+2 {
				t.Errorf("Expected row %d was %d", r+2, f[r].Row)
				return
			}
		}
	}
}

func Test_BasicHorizontal(t *testing.T) {
	for i := 0; i+4 < model.Column; i++ {
		g := new(model.Game)
		for c := i; c < i+4; c++ {
			g.DropPiece(i%2 == 0, c)
		}

		f := horizontal4InRow(g.GetBoard())
		if f == nil {
			t.Error("Should detect horizontal four in a row")
			return
		}

		for l := range 4 {
			if f[l].Row != 0 {
				t.Errorf("Expected row 0 was %d", f[l].Row)
				return
			}

			if f[l].Column != l+i {
				t.Errorf("Expected column %d was %d", l+i, f[l].Column)
				return
			}
		}
	}
}
