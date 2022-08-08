package main

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket"
)

func serveWs(pool *websocket.ChatPool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chat WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := websocket.NewClient(r, conn, pool)

	pool.Register <- client

	client.Read()
}

func serveGameWs(pool *websocket.GamePool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Game WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := websocket.NewGameClient(r, conn, pool)

	pool.Register <- client

	client.Read()
}

func setupRoutes() {
	chatPool := websocket.NewChatPool()
	go chatPool.Start()

	gamePool := websocket.NewGamePool()
	go gamePool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chatPool, w, r)
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
