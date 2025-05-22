package gameflow

import (
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
