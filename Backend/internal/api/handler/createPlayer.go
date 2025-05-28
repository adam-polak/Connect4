package handler

import (
	"io"
	"net/http"
	"strings"
)

type CreatePlayerHandler struct{}

type player struct {
	loginKey string
	username string
}

func (c *CreatePlayerHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if strings.Compare(r.Method, "POST") != 0 {
		wr.WriteHeader(400)
		return
	}

	p := getPlayer(r.Body)
	if p == nil {
		wr.WriteHeader(400)
		return
	}
}

func getPlayer(b io.ReadCloser) *player {
	return nil
}
