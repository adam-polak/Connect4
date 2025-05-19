package logic

import (
	"connect4/server/internal/game/model"
	"testing"
)

func Test_BasicVertical(t *testing.T) {
	for c := 0; c < model.Column; c++ {
		g := new(model.Game)
		for range target {
			g.DropPiece(c%2 == 0, c)
		}

		f := vertical4InRow(g.GetBoard())
		if f == nil {
			t.Error("Should detect vertical four in a row")
			return
		}

		for r := range target {
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
			if r < model.Row-target {
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

		for r := range target {
			if f[r].Column != c {
				t.Errorf("Expected column %d was %d", c, f[r].Column)
				return
			}

			if f[r].Row != r+model.Row-target {
				t.Errorf("Expected row %d was %d", r+2, f[r].Row)
				return
			}
		}
	}
}

func Test_BasicHorizontal(t *testing.T) {
	for i := 0; i+target-1 < model.Column; i++ {
		g := new(model.Game)
		for c := i; c < i+target; c++ {
			g.DropPiece(i%2 == 0, c)
		}

		f := horizontal4InRow(g.GetBoard())
		if f == nil {
			t.Error("Should detect horizontal four in a row")
			return
		}

		for l := range target {
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

func Test_LeftToRightDiagonal(t *testing.T) {
	for c := range model.Column - target + 1 {
		g := new(model.Game)
		for r := range target {
			for i := range r {
				g.DropPiece(i+r%2 == 0, c+r)
			}

			g.DropPiece(c%2 == 0, c+r)
		}

		f := diagonal4InRow(g.GetBoard())
		if f == nil {
			t.Error("Should detect diagonal four in a row")
			return
		}

		for i := range target {
			if f[i].Column != c+i {
				t.Errorf("Expected column %d was %d", c, f[i].Column)
				return
			}

			if f[i].Row != i {
				t.Errorf("Expected row %d was %d", i, f[i].Row)
				return
			}
		}
	}
}
