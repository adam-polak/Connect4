package player

import (
	"errors"
	"strings"
)

type player struct {
	uid string
	key string
}

type playerNode struct {
	player player
	left   *playerNode
	right  *playerNode
}

var players *playerNode = nil

func NewPlayer(loginKey string) {
	// add player to tree
}

func getPlayer(p *playerNode, key string) *string {
	if p == nil {
		return nil
	}

	cmp := strings.Compare(p.player.key, key)
	if cmp < 0 {
		return getPlayer(p.left, key)
	} else if cmp > 0 {
		return getPlayer(p.right, key)
	}

	return &p.player.uid
}

func GetPlayer(loginKey string) (string, error) {
	if players == nil {
		return "", errors.New("player list must be initialized")
	}

	uid := getPlayer(players, loginKey)
	if uid == nil {
		return "", errors.New("player doesn't exist")
	}

	return *uid, nil
}
