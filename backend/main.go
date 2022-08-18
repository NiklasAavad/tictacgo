package main

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/chat"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/gamesocket"
)

func serveWs(pool websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := pool.NewClient(r, conn)
	pool.Register(client)

	client.Read()
}

func startPools() (*chat.ChatPool, *gamesocket.GamePool) {
	chatPool := chat.NewChatPool()
	gamePool := gamesocket.NewGamePool()

	go chatPool.Start()
	go gamePool.Start()

	return chatPool, gamePool
}

func setupRoutes() {
	chatPool, gamePool := startPools()

	http.HandleFunc("/chatws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chatPool, w, r)
	})

	http.HandleFunc("/gamews", func(w http.ResponseWriter, r *http.Request) {
		serveWs(gamePool, w, r)
	})
}

func main() {
	fmt.Println("Go App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
