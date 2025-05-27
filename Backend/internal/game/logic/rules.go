package logic

import (
	"connect4/server/internal/game/model"
	"errors"
)

const (
	target = 4
)

type board [model.Column][model.Row]uint8

type Location struct {
	Column int
	Row    int
}

type FourInARow [target]Location

func diagonal4InRow(b board) *FourInARow {
	l := leftDiagonal4InRow(b)
	if l != nil {
		return l
	}

	return rightDiagonal4InRow(b)
}

func leftDiagonal4InRow(b board) *FourInARow {
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
			for p == b[loc.Column][loc.Row] && loc.Column < model.Column && loc.Row < model.Row {
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
				loc.Row++
				if len == target {
					return f
				}
			}
		}
	}

	return nil
}

func rightDiagonal4InRow(b board) *FourInARow {
	f := new(FourInARow)
	for c := model.Column - 1; c >= 0; c++ {
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
			for p == b[loc.Column][loc.Row] && loc.Column < model.Column && loc.Row < model.Row {
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
				loc.Column--
				loc.Row++
				if len == target {
					return f
				}
			}
		}
	}

	return nil
}

func vertical4InRow(b board) *FourInARow {
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

func horizontal4InRow(b board) *FourInARow {
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

func Has4InARow(g model.Game) *FourInARow {
	b := g.GetBoard()
	// diag := diagonal4InRow(b)
	// if diag != nil {
	// 	return diag
	// }

	// vert := vertical4InRow(b)
	// if vert != nil {
	// 	return vert
	// }

	// hrz := horizontal4InRow(b)
	// if hrz != nil {
	// 	return hrz
	// }

	// return nil

	return has4InARow(b)
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

func upDirection() direction {
	return direction{
		dx: 0,
		dy: 1,
	}
}

func rightDirection() direction {
	return direction{
		dx: 1,
		dy: 0,
	}
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
		dy: 1,
	}
}

func checkDirection(b board, f *FourInARow, col int, row int, dir direction) bool {
	if b[col][row] == model.NoPlayer {
		return false
	}

	f[0] = Location{
		Column: col,
		Row:    row,
	}

	found, _ := targetSearch(b, f, 1, dir)
	return found
}

func has4InARow(b board) *FourInARow {
	directions := []direction{rightDirection(), upRightDirection(), downRightDirection(), upDirection()}
	f := new(FourInARow)
	for c := range model.Column {
		canRight, canUp := c+target-1 < model.Column, true
		for r := range model.Row {
			canUp = r+target-1 < model.Row
			if !canRight && !canUp {
				break
			}

			if canRight && canUp {
				for i := range directions {
					if checkDirection(b, f, c, r, directions[i]) {
						return f
					}
				}
			} else if canRight && checkDirection(b, f, c, r, directions[0]) {
				return f
			} else if canUp && checkDirection(b, f, c, r, directions[3]) {
				return f
			}
		}
	}

	return nil
}

// Search from the current index on the board and travel in the specified direction,
// if the target length is achieved return true
//
// - cur index must start at 1, the 0th index should be set before the function is entered
func targetSearch(b board, f *FourInARow, cur int, dir direction) (bool, error) {
	// set cur
	f[cur].Column = f[cur-1].Column + dir.dx
	f[cur].Row = f[cur-1].Row + dir.dy

	// check if location is in range
	if f[cur].Column < 0 || f[cur].Column >= model.Column || f[cur].Row < 0 || f[cur].Row >= model.Row {
		return false, errors.New("not in range")
	}

	// check if last position is equal
	if b[f[cur].Column][f[cur].Row] != b[f[cur-1].Column][f[cur-1].Row] {
		return false, nil
	}

	// check if target reached
	if cur == target-1 {
		return true, nil
	}

	return targetSearch(b, f, cur+1, dir)
}
