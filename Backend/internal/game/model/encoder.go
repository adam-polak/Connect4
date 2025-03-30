package model

const (
	ColumnShift = 9
)

// Encode the game into a 64 bit integer.
// - Each of the 7 columns is split into 9 bits
// - The first 3 bits signify the position of the last piece
// in the column
//   - 000 ==> no pieces in the column
//   - 110 ==> last piece is in row 5
//
// - The last 6 bits signify which players piece is where
//   - 0 for player 1 and 1 for player 2
func (g Game) Encode() uint64 {
	var encoding uint64 = 0

	for i := range Column {
		var shift, last = ColumnShift * i, Row
		var v uint64 = 0 // v = encoded value of column
		for j := range Row {
			if g.board[i][j] == 0 {
				last = j
				break
			} else if g.board[i][j] == 2 {
				v = v | 1<<(shift+j)
			}
		}

		v = v | (uint64(last) << (shift + Row))
		encoding = encoding | v
	}

	return encoding
}
