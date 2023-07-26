package main

import (
	"fmt"
	"net/http"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/chat"
	"github.com/NiklasPrograms/tictacgo/backend/pkg/websocket/gamesocket"
)

func startChatPool() *chat.ChatPool {
	chatPool := chat.NewChatPool()
	go chatPool.Start()
	return chatPool
}

func startGamePool() *gamesocket.GamePool {
	channelStrategy := gamesocket.NewConcurrentChannelStrategy()
	gamePool := gamesocket.NewGamePool(channelStrategy)
	go gamePool.Start()
	return gamePool
}

func setupRoutes() {
	chatPool := startChatPool()
	http.HandleFunc("/chatws", func(w http.ResponseWriter, r *http.Request) {
		chatPool.ServeWs(w, r)
	})

	gamePool := startGamePool()
	http.HandleFunc("/gamews", func(w http.ResponseWriter, r *http.Request) {
		gamePool.ServeWs(w, r)
	})
}

func main() {
	fmt.Println("Go App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
