package handler

import (
	"connect4/server/internal/game/gameflow"
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

	_, err := gameflow.GetPlayer(key)
	if err != nil {
		wr.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP --> WebSocket
	conn, err := upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, buf, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(messageType, buf); err != nil {
			log.Println(err)
			return
		}
	}
}
