package model

func Decode(n uint64) *Game {
	g := new(Game)

	for i := range Column {
		col := (n >> (ColumnShift * i)) & ((1 << ColumnShift) - 1)
		full := (col >> Row) & ((1 << 3) - 1)
		for j := range full {
			player := (col >> (full - j - 1)) & 1

			g.DropPiece(player == 0, i)
		}
	}

	return g
}
