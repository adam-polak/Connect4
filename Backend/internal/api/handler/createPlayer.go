package handler

import (
	"connect4/server/internal/game/gameflow"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type player struct {
	LoginKey string
	Username string
}

func CreatePlayerHandler(wr http.ResponseWriter, r *http.Request) {
	if strings.Compare(r.Method, "POST") != 0 {
		wr.WriteHeader(400)
		return
	}

	p := getPlayer(r.Body)
	if p == nil {
		wr.WriteHeader(400)
		return
	}

	// Create new player
	_, err := gameflow.NewPlayer(p.LoginKey, p.Username)
	if err != nil {
		// player with key already exists with different username
		wr.WriteHeader(http.StatusConflict)
		return
	}

	wr.WriteHeader(http.StatusOK)
}

func getPlayer(b io.ReadCloser) *player {
	defer b.Close()
	s, err := io.ReadAll(b)
	if err != nil {
		return nil
	}

	var p player
	err = json.Unmarshal(s, &p)
	if err != nil {
		return nil
	}

	return &p
}
