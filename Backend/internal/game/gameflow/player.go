package gameflow

import (
	"errors"
	"fmt"
	"strings"
)

const (
	PlayersListNotInitialized = "player list not initialized"
	PlayerDoesNotExist        = "player doesn't exist"
	PlayerAlreadyExists       = "player already exists"
)

type Player struct {
	Username string
	game     *GameOrchestrator
}

type playerNode struct {
	key    string
	player *Player
	left   *playerNode
	right  *playerNode
}

var players *playerNode = nil

func NewPlayer(key string, username string) (*Player, error) {
	p := getPlayer(key)
	if p != nil && strings.Compare(p.Username, username) != 0 {
		return nil, errors.New(PlayerAlreadyExists)
	} else if p != nil {
		return p, nil
	}

	p = &Player{
		Username: username,
	}

	addPlayer(key, p)

	return p, nil
}

func addPlayer(key string, p *Player) {
	pn := &playerNode{
		key:    key,
		player: p,
	}
	n := players
	if n == nil {
		players = pn
		return
	}

	for {
		cmp := strings.Compare(n.key, key)
		if cmp < 0 {
			if n.left == nil {
				n.left = pn
				return
			}

			n = n.left
		} else if cmp > 0 {
			if n.right == nil {
				n.right = pn
				return
			}

			n = n.right
		} else {
			panic("key already exists in players list")
		}
	}
}

func getPlayerRecursive(p *playerNode, key string) *Player {
	if p == nil {
		return nil
	}

	cmp := strings.Compare(p.key, key)
	if cmp < 0 {
		return getPlayerRecursive(p.left, key)
	} else if cmp > 0 {
		return getPlayerRecursive(p.right, key)
	}

	return p.player
}

func getPlayer(key string) *Player {
	return getPlayerRecursive(players, key)
}

func GetPlayer(key string) (*Player, error) {
	if players == nil {
		return nil, errors.New(PlayersListNotInitialized)
	}

	p := getPlayer(key)
	if p == nil {
		return nil, errors.New(PlayerDoesNotExist)
	}

	return p, nil
}

func replaceNode(n *playerNode) *playerNode {
	if n == nil {
		return nil
	} else if n.right == nil {
		return n.left
	}

	r := n.right
	for r.left != nil {
		prev := r
		r = r.left
		if r.left == nil {
			prev.left = nil
		}
	}

	r.left = n.left

	return r
}

func doRemove(n *playerNode, key string) bool {
	if n == nil {
		return false
	}

	cmp := strings.Compare(n.key, key)
	if cmp == 0 {
		if n == players {
			players = replaceNode(n)
			fmt.Println(players.key)
		}

		return true
	} else if cmp < 0 && doRemove(n.left, key) {
		n.left = replaceNode(n.left)
	} else if cmp > 0 && doRemove(n.right, key) {
		n.right = replaceNode(n.right)
	}

	return false
}

func RemovePlayer(key string) error {
	p, _ := GetPlayer(key)
	if p == nil {
		return errors.New(PlayerDoesNotExist)
	}

	doRemove(players, key)

	return nil
}

func printPlayers(n *playerNode, level int) {
	if n == nil {
		return
	}

	fmt.Printf("%d: %s\n", level, n.key)
	printPlayers(n.left, level+1)
	printPlayers(n.right, level+1)
}

func PrintPlayers() {
	fmt.Println("Players\n-------------------------")
	printPlayers(players, 1)
	fmt.Println("-------------------------")
}

func (p *Player) SetGame(g *GameOrchestrator) {
	p.game = g
}
