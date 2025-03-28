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
func (g Game) encode() uint64 {
	var encoding uint64 = 0

	for i := range Column {
		var shift, last uint32 = ColumnShift * (uint32(i) + 1), 0
		var v uint64 = 0
		for j := range Row {
			if g.board[i][j] == 0 {
				last = uint32(j)
				break
			} else if g.board[i][j] == 2 {
				v = 1 << (shift + uint32(j))
			}
		}

		v = v | (uint64(last) << (shift + Row))
		encoding = encoding | v
	}

	return encoding
}
