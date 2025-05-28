package entry

import (
	"connect4/server/internal/api/handler"
	"log"
	"net/http"
)

func StartServer() {
	log.Println("🦍 Starting server")
	// Add handlers
	http.Handle("/game", new(handler.GameHandler))
	// Listen to port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
