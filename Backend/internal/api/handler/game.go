package handler

import (
	"connect4/server/internal/game/gameflow"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func writePump(c *websocket.Conn, ch chan string) {

}

func readPump(c *websocket.Conn, ch chan string) {

}

func observePlayer(p gameflow.Player, ch chan string) {
	playerObserver := func(action interface{}) {
		switch v := action.(type) {
		case gameflow.GameReady:
			ch <- "game is ready to play"
		case gameflow.GameOver:
			ch <- "game is over"
			fmt.Println("game is over")
		case gameflow.OpponentPlayed:
			ch <- fmt.Sprintf("opponent played in column %d", v.Column)
		default:
			ch <- "case not recognized"
		}
	}

	p.RegisterObserver(playerObserver)
}

func GameHandler(wr http.ResponseWriter, r *http.Request) {
	enableCors(wr)

	if len(r.Method) != 0 && strings.Compare(r.Method, "GET") != 0 {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("key")
	if len(key) == 0 {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := gameflow.GetPlayer(key)
	if err != nil {
		wr.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Join a game
	gameflow.JoinGame(p)

	// Upgrade HTTP --> WebSocket
	conn, err := upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	writeChannel := make(chan string)

	observePlayer(*p, writeChannel)
	go writePump(conn, writeChannel)
	go readPump(conn, writeChannel)
}
