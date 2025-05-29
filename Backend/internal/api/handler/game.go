package handler

import (
	"connect4/server/internal/game/gameflow"
	"encoding/json"
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

type sendError struct {
	Error string
}

type doPlay struct {
	Column int `json:"column"`
}

type logPlay struct {
	Column int  `json:"column"`
	IsSelf bool `json:"isSelf"` // boolean to tell client which player made the change
}

type playerJoined struct {
	Username string `json:"username"`
}

type informPlayer struct {
	SelfUsername     string `json:"selfUsername"`
	OpponentUsername string `json:"opponentUsername"`
}

func writePump(c *websocket.Conn, ch chan []byte) {
	defer c.Close()

	for {
		msg, ok := <-ch
		if !ok {
			c.WriteJSON(sendError{Error: "A server error occurred"})
			return
		}

		c.WriteMessage(websocket.TextMessage, msg)
	}
}

func readPump(c *websocket.Conn, p *gameflow.Player, ch chan []byte) {
	defer c.Close()

	for {
		var j interface{}
		err := c.ReadJSON(&j)
		if err != nil {
			log.Println(err)
			return
		}

		switch v := j.(type) {
		case doPlay:
			err = p.PlayPiece(v.Column)
			if err != nil {
				if strings.Compare(err.Error(), gameflow.FailedAction) != 0 {
					log.Println(err)
					return
				}

				b, err := json.Marshal(sendError{Error: gameflow.FailedAction})
				if err != nil {
					log.Println(err)
					return
				}

				ch <- b
				continue
			}

			b, err := json.Marshal(logPlay{Column: v.Column, IsSelf: true})
			if err != nil {
				log.Println(err)
				return
			}

			ch <- b
		default:
			log.Println("JSON not recognized")
			return
		}
	}
}

func getPlayerObserver(ch chan []byte) func(interface{}) {
	return func(action interface{}) {
		switch v := action.(type) {
		case gameflow.GameReady:
			b, err := json.Marshal(playerJoined{Username: v.OpponentUsername})
			if err != nil {
				log.Println(err)
				return
			}

			ch <- b
		case gameflow.GameOver:
			ch <- []byte("game is over")
		case gameflow.OpponentPlayed:
			b, err := json.Marshal(logPlay{Column: v.Column, IsSelf: false})
			if err != nil {
				log.Println(err)
				return
			}

			ch <- b
		default:
			ch <- []byte("case not recognized")
		}
	}
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

	ch := make(chan []byte)
	obs := getPlayerObserver(ch)
	p, err := gameflow.GetPlayer(key, &obs)
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

	go writePump(conn, ch)
	go readPump(conn, p, ch)

	opp := p.GetOpponentUsername()
	var b []byte
	if opp == nil {
		b, err = json.Marshal(informPlayer{
			SelfUsername:     p.Username,
			OpponentUsername: "",
		})
	} else {
		b, err = json.Marshal(informPlayer{
			SelfUsername:     p.Username,
			OpponentUsername: *opp,
		})
	}

	if err != nil {
		log.Println(err)
		return
	} else {
		ch <- b
	}

	// Join a game
	gameflow.JoinGame(p)
}
