package handler

import (
	"connect4/server/internal/game/gameflow"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var playerLock sync.Mutex

type sendError struct {
	Error string `json:"error"`
}

type doPlay struct {
	Column int `json:"column"`
}

type newGame struct {
	Requested bool `json:"requested"`
}

type logPlay struct {
	Column int  `json:"column"`
	IsSelf bool `json:"isSelf"` // boolean to tell client which player made the change
}

type playerJoined struct {
	Username string `json:"username"`
}

type gameOver struct {
	YouWin bool `json:"youWin"`
}

type informPlayer struct {
	SelfUsername     string `json:"selfUsername"`
	OpponentUsername string `json:"opponentUsername"`
	GameWinner       int    `json:"gameWinner"`
	Starter          bool   `json:"starter"`
	YourTurn         bool   `json:"yourTurn"`
	Plays            []int  `json:"plays"`
}

func writePump(c *websocket.Conn, msgs chan []byte, username string) {
	ticker := time.NewTicker(54 * time.Second)

	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				c.WriteJSON(sendError{Error: "A server error occurred"})
				return
			}

			c.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Write error for %s: %v", username, err)
				return
			}
		case <-ticker.C:
			c.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readPump(c *websocket.Conn, p *gameflow.Player, ch chan []byte) {
	defer c.Close()

	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}

		log.Printf("Incoming: %s\n", string(msg))

		var play doPlay
		decoder := json.NewDecoder(strings.NewReader(string(msg)))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&play)
		if err == nil {
			playerLock.Lock()
			err = p.PlayPiece(play.Column)
			playerLock.Unlock()
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

			b, err := json.Marshal(logPlay{Column: play.Column, IsSelf: true})
			if err != nil {
				log.Println(err)
				return
			}

			ch <- b
			continue
		}

		var newGame newGame
		err = json.Unmarshal(msg, &newGame)
		if err == nil {
			playerLock.Lock()
			p.FindNewGame()
			playerLock.Unlock()

			continue
		}
	}
}

func getPlayerObserver(ch chan []byte, playerKey string) func(interface{}) {
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
			b, err := json.Marshal(gameOver{YouWin: strings.Compare(playerKey, v.Winner) == 0})
			if err != nil {
				log.Println(err)
				return
			}

			ch <- b
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
	obs := getPlayerObserver(ch, key)
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

	go writePump(conn, ch, p.Username)
	go readPump(conn, p, ch)

	playerLock.Lock()
	opp := p.GetOpponentUsername()
	playerLock.Unlock()
	var b []byte
	if opp == nil {
		b, err = json.Marshal(informPlayer{
			SelfUsername:     p.Username,
			OpponentUsername: "",
			GameWinner:       0,
			Plays:            []int{},
		})
	} else {
		playerLock.Lock()
		w := p.GetGameWinner()
		playerLock.Unlock()
		var gw int
		if len(w) == 0 {
			// no winner
			gw = 0
		} else if strings.Compare(w, key) == 0 {
			// winner is player 1 i.e. self
			gw = 1
		} else {
			// winner is player 2 i.e. opponent
			gw = 2
		}

		playerLock.Lock()
		s := p.GetStartingPlayer()
		st := strings.Compare(s, key) == 0

		b, err = json.Marshal(informPlayer{
			SelfUsername:     p.Username,
			OpponentUsername: *opp,
			GameWinner:       gw,
			Starter:          st,
			YourTurn:         p.IsYourTurn(),
			Plays:            p.GetPlays(),
		})
		playerLock.Unlock()
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
