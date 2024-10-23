package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Type    string
	Content interface{}
}

var connected_players []Player
var played_cards map[Card]Player

func connectPlayer(conn *websocket.Conn) {
	player := &Player{
		strings.Split(conn.RemoteAddr().String(), ":")[0],
		conn,
		make([]Card, 5),
		false,
	}
	connected_players = append(connected_players, *player)
	message := &Message{
		"PlayerJoined",
		fmt.Sprintf("Player with Addr %s has joined.", player.Addr),
	}
	fmt.Println("Player with Addr " + player.Addr + " has joined.")
	broadcastMessage(*message)
}

func disconnectPlayer(conn *websocket.Conn) {
	for i, player := range connected_players {
		if player.Conn == conn {
			connected_players = append(connected_players[:i], connected_players[i+1:]...)
			message := &Message{
				"PlayerLeft",
				fmt.Sprintf("Player with Addr %s has left.", player.Addr),
			}
			broadcastMessage(*message)
			break
		}
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	connectPlayer(ws)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			disconnectPlayer(conn)
			fmt.Println("Error reading message:", err)
			return
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println("Error unmarshaling message:", err)
			continue
		}
	}
}

func broadcastMessage(message Message) {
	for _, player := range connected_players {
		msg, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		player.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func checkIfAllPlayed() bool {
	return len(played_cards) == len(connected_players)
}
