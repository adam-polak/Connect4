package model

import (
	"encoding/binary"
	"strings"
	"testing"
)

func ByteString(n uint64) string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	var sb strings.Builder
	for i := 0; i < len(b); i++ {
		for j := 7; j >= 0; j-- {
			bit := b[i] & (1 << j)
			if bit > 0 {
				sb.WriteString("1")
			} else {
				sb.WriteString("0")
			}
		}
		sb.WriteString(" ")
	}

	return sb.String()
}

// TestGame_EncodeEmpty tests that the encoding is accurate
// for an empty Connect4 board
func TestGame_EncodeEmpty(t *testing.T) {
	g := new(Game)
	encoded := g.Encode()

	if encoded != 0 {
		t.Errorf("Expected new board encoded to equal 0, was %d", encoded)
	}
}

// TestGame_EncodeSingleColumnFull tests that encoding is accurate for
// a single column being full
func TestGame_EncodeSingleColumnFull(t *testing.T) {
	var fullVal uint64 = Row // integer indicating column is full

	for i := range Column {
		g := new(Game)
		for j := range Row {
			g.board[i][j] = 1 // manually fill the ith column
		}

		encoded := g.Encode()
		shift := ColumnShift * i
		var expected uint64 = fullVal << (shift + Row)
		if encoded != expected {
			t.Errorf(
				"\n\nExpected:\n\t %s \nActual:\n\t %s",
				ByteString(expected),
				ByteString(encoded),
			)
		}
	}
}
