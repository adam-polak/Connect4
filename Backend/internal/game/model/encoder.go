package model

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
	// TODO
	return 0
}
