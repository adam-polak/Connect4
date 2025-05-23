package gameflow

import (
	"fmt"
	"testing"
)

func cleanupPlayers() {
	players = nil
}

func Test_NewPlayerBasic(t *testing.T) {
	defer cleanupPlayers()
	p, err := NewPlayer("hello", "bob")
	if p == nil || err != nil {
		t.Error("Should create a player")
		return
	}
}

func Test_NewPlayerFailsWithDifferentUsername(t *testing.T) {
	defer cleanupPlayers()
	k := "key"
	NewPlayer(k, "bob")
	_, err := NewPlayer(k, "sue")
	if err == nil {
		t.Error("Creating a player with the same key and different user id should fail")
		return
	}
}

func Test_NewPlayerIsTheSameObject(t *testing.T) {
	defer cleanupPlayers()
	k := "superman"
	u := "bob"
	p1, e1 := NewPlayer(k, u)
	p2, e2 := NewPlayer(k, u)
	if p1 == nil || p2 == nil || e1 != nil || e2 != nil {
		t.Error("Shouldn't fail to create same player")
		return
	}

	if p1 != p2 {
		t.Error("Should be same object returned when creating same player")
		return
	}
}

func Test_NewPlayerEqualsGetPlayer(t *testing.T) {
	defer cleanupPlayers()
	k := "supersecret"
	u := "bob"
	p1, e1 := NewPlayer(k, u)
	if p1 == nil || e1 != nil {
		t.Error("Shouldn't fail to create a player")
		return
	}

	p2, e2 := GetPlayer(k)
	if p2 == nil || e2 != nil {
		t.Error("Shouldn't fail to get player")
		return
	}

	if p1 != p2 {
		t.Error("Getting player after creating new player should be equal")
		return
	}
}

func Test_RemovePlayerShouldFailWhenEmpty(t *testing.T) {
	if RemovePlayer("test") == nil {
		t.Error("Remove player should throw")
	}
}

func Test_RemovePlayerShouldFailWhenNoPlayer(t *testing.T) {
	defer cleanupPlayers()
	NewPlayer("key", "bob")
	if RemovePlayer("test") == nil {
		t.Error("Remove player should throw")
	}
}

func Test_RemovePlayerSmall(t *testing.T) {
	defer cleanupPlayers()
	arr := [...]string{"i", "c", "a", "d", "e", "u", "w", "o"}
	for i := range arr {
		p, err := NewPlayer(arr[i], arr[i])
		if p == nil || err != nil {
			t.Error("Failed to create player")
			return
		}
	}

	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		if RemovePlayer(arr[i]) != nil || RemovePlayer(arr[j]) != nil {
			t.Error("Failed to remove player")
			return
		}

		for x, y := i+1, j-1; x < y; x, y = x+1, y-1 {
			p1, _ := GetPlayer(arr[x])
			if p1 == nil {
				t.Error("Player should exist")
				return
			}

			p2, _ := GetPlayer(arr[y])
			if p2 == nil {
				t.Error("Player should exist")
				return
			}
		}
	}
}

func Test_RemovePlayerLarge(t *testing.T) {
	defer cleanupPlayers()
	arr := [...]string{"i", "c", "a", "d", "e", "u", "w", "o"}
	len := len(arr)
	size := 5
	for i := range size {
		s := ""
		for j := range (i / len) + 1 {
			s += arr[(i+j)%len]
		}

		p, err := NewPlayer(s, "joe")
		if p == nil || err != nil {
			t.Error("Failed to create player")
			return
		}
	}

	PrintPlayers()

	for i := range size {
		s := ""
		for j := range (i / len) + 1 {
			s += arr[(i+j)%len]
		}

		fmt.Printf("Removing %s\n", s)
		fmt.Println("----------------")

		if RemovePlayer(s) != nil {
			t.Error("Failed to remove player")
			return
		}

		PrintPlayers()

		for j := i + 1; j < size; j++ {
			s2 := ""
			for x := range (j / len) + 1 {
				s2 += arr[(j+x)%len]
			}

			p, _ := GetPlayer(s2)
			if p == nil {
				t.Error("Player should exist")
				return
			}
		}
	}
}
