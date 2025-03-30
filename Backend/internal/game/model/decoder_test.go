package model

import (
	"math/rand"
	"testing"
)

func (g Game) equals(g2 *Game) bool {
	for i := range Column {
		for j := range Row {
			if g.board[i][j] != g2.board[i][j] {
				return false
			}
		}
	}

	return true
}

func TestGame_DecodeEmpty(t *testing.T) {
	g := new(Game)

	if !g.equals(Decode(0)) {
		t.Error("Decoding 0 should equal a new game")
	}
}

func TestGame_EncodeAndDecode(t *testing.T) {
	g := new(Game)

	for i := range Column {
		for j := range rand.Int31n(Row) {
			g.DropPiece(i+int(j)%2 == 0, i)
		}
	}

	gEncode := g.Encode()
	decoded := Decode(gEncode)
	dEncode := decoded.Encode()
	if !g.equals(decoded) {
		t.Errorf(
			"Decoding a game encoding should equal the same game\nExpected:\n\t%s\nActual:\n\t%s\n",
			ByteString(gEncode),
			ByteString(dEncode),
		)
	}
}
