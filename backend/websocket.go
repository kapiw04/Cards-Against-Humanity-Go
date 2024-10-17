package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Type    string
	Content string
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			panic(err)
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			panic(err)
		}

		switch message.Type {
		case "Test":
			fmt.Println("[TEST] Received test message")
			fmt.Println(message.Content)
		}
	}
}
