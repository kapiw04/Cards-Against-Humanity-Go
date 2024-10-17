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

func connectPlayer(conn *websocket.Conn) {
	ID := strings.Split(conn.RemoteAddr().String(), ":")[1]
	player := &Player{
		ID,
		conn,
		make([]Card, 5),
	}
	connected_players = append(connected_players, *player)
	message := &Message{
		"PlayerJoined",
		fmt.Sprintf("Player with ID %s has joined.", player.ID),
	}
	broadcastMessage(*message)
}

func disconnectPlayer(conn *websocket.Conn) {
	for i, player := range connected_players {
		if player.Conn == conn {
			connected_players = append(connected_players[:i], connected_players[i+1:]...)
			message := &Message{
				"PlayerLeft",
				fmt.Sprintf("Player with ID %s has left.", player.ID),
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
			break
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			panic(err)
		}

		switch message.Type {
		case "GetHand":
			for _, player := range connected_players {
				if player.Conn == conn {
					writePlayersHand(player)
					break
				}
			}
		}

	}
}

func writePlayersHand(player Player) {
	handJSON := HandJSON{
		Cards: make([]CardJSON, len(player.Hand)),
	}
	for i, card := range player.Hand {
		fmt.Println(card)
		handJSON.Cards[i] = CardJSON(card)
	}

	message := &Message{
		Type:    "SendPlayerHand",
		Content: handJSON,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	err = player.Conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		panic(err)
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
