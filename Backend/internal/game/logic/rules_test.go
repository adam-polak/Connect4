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

		f := Has4InARow(g.GetBoard())
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

		f := Has4InARow(g.GetBoard())
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
				t.Errorf("Expected row %d was %d", r+model.Row-target, f[r].Row)
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

		f := Has4InARow(g.GetBoard())
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

func Test_UpRightDiagonal(t *testing.T) {
	for c := 0; c+target-1 < model.Column; c++ {
		for r := 0; r+target-1 < model.Row; r++ {
			loc := Location{
				Column: c,
				Row:    r,
			}

			b := new(Board)
			dir := upRightDirection()
			for range target {
				b[loc.Column][loc.Row] = model.PlayerOne

				loc.Column = loc.Column + dir.dx
				loc.Row = loc.Row + dir.dy
			}

			f := Has4InARow(*b)
			if f == nil {
				t.Error("Should detect diagonal four in a row")
				return
			}

			for i := range target {
				if f[i].Column != loc.Column-target+i {
					t.Errorf("Expected column %d was %d", loc.Column-target+i, f[i].Column)
					return
				}

				if f[i].Row != loc.Row-target+i {
					t.Errorf("Expected row %d was %d", loc.Row-target+i, f[i].Row)
					return
				}
			}
		}
	}
}

func Test_DownRightDiagonal(t *testing.T) {
	for c := 0; c+target-1 < model.Column; c++ {
		for r := model.Row - 1; r-target+1 >= 0; r-- {
			loc := Location{
				Column: c,
				Row:    r,
			}

			b := new(Board)
			dir := downRightDirection()
			for range target {
				b[loc.Column][loc.Row] = model.PlayerOne
				loc.Column = loc.Column + dir.dx
				loc.Row = loc.Row + dir.dy
			}

			f := Has4InARow(*b)
			if f == nil {
				t.Error("Should detect diagonal four in a row")
				return
			}

			for i := range target {
				if f[i].Column != loc.Column-target+i {
					t.Errorf("Expected column %d was %d", loc.Column-target+i, f[i].Column)
					return
				}

				if f[i].Row != loc.Row+target-i {
					t.Errorf("Expected row %d was %d", loc.Row-i, f[i].Row)
					return
				}
			}
		}
	}
}

func Test_Efficiency(t *testing.T) {
	for range 100000 {
		Test_BasicVertical(t)
		Test_BasicHorizontal(t)
		Test_Last4Vertical(t)
		Test_UpRightDiagonal(t)
		Test_DownRightDiagonal(t)
	}
}
