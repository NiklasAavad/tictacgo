package main

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/chat"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/gamesocket"
)

func serveChatWs(pool *chat.ChatPool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chat WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := chat.NewClient(r, conn, pool)

	pool.Register <- client

	client.Read()
}

func serveGameWs(pool websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Game WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := pool.NewClient(r, conn)
	pool.Register(client)

	client.Read()
}

func setupRoutes() {
	chatPool := chat.NewChatPool()
	go chatPool.Start()

	gamePool := gamesocket.NewGamePool()
	go gamePool.Start()

	http.HandleFunc("/chatws", func(w http.ResponseWriter, r *http.Request) {
		serveChatWs(chatPool, w, r)
	})

	http.HandleFunc("/gamews", func(w http.ResponseWriter, r *http.Request) {
		serveGameWs(gamePool, w, r)
	})
}

func main() {
	fmt.Println("Go App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
