package gameflow

import (
	"errors"
	"strings"
	"sync"
)

const (
	PlayerDoesNotExist  = "player doesn't exist"
	PlayerAlreadyExists = "player already exists"
)

type Player struct {
	Username string
	observer *func(interface{})
	game     *GameOrchestrator
}

type playerNode struct {
	key    string
	player *Player
	left   *playerNode
	right  *playerNode
}

var players *playerNode = nil

var playersLock sync.Mutex

func NewPlayer(key string, username string) (*Player, error) {
	playersLock.Lock()
	defer playersLock.Unlock()

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

func GetPlayer(key string, observer *func(interface{})) (*Player, error) {
	p := getPlayer(key)
	if p == nil {
		return nil, errors.New(PlayerDoesNotExist)
	}

	p.observer = observer

	return p, nil
}

func replaceNode(n *playerNode) *playerNode {
	if n == nil {
		return nil
	} else if n.right == nil {
		return n.left
	}

	prev := n
	r := n.right
	for r.left != nil {
		prev = r
		r = r.left
	}

	// replace left branch of r with node's left branch
	r.left = n.left
	// replace previous node left with right branch of r
	prev.left = r.right
	// replace right branch of r with node's right branch
	if n.right != r {
		r.right = n.right
	}

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
	playersLock.Lock()
	defer playersLock.Unlock()

	p := getPlayer(key)
	if p == nil {
		return errors.New(PlayerDoesNotExist)
	}

	doRemove(players, key)

	return nil
}

func removeAllPlayers() {
	players = nil
}
