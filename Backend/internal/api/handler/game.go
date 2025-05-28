package handler

import (
	"net/http"
	"strings"
)

type GameHandler struct{}

func (g *GameHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if len(r.Method) != 0 && strings.Compare(r.Method, "GET") != 0 {
		wr.WriteHeader(400)
		return
	}

}
