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
var played_cards []CardPlayedJSON

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
		case "GetBlackCard":
			writeBlackCard()

		case "PlayCard":
			handleCardPlayed(conn, message)

		case "GetAllCards":
			writeAllPlayedCards()
		}
	}
}

func writePlayersHand(player Player) {
	handJSON := HandJSON{
		Cards: make([]CardJSON, len(player.Hand)),
	}
	for i, card := range player.Hand {
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

func writeBlackCard() {
	cardJSON := CardJSON(currentBlackCard)

	message := &Message{
		Type:    "SendBlackCard",
		Content: cardJSON,
	}

	broadcastMessage(*message)
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

func handleCardPlayed(conn *websocket.Conn, message Message) {
	var card Card
	println(message.Content)
	content, ok := message.Content.([]byte)
	if !ok {
		panic("Error converting content to JSON")
	}
	err := json.Unmarshal(content, &card)
	if err != nil {
		panic(err)
	}
	cardPlayed := &CardPlayedJSON{
		Text:       card.Text,
		OwnersConn: conn,
	}
	for i, player := range connected_players {
		if player.Conn == conn {
			connected_players[i].Hand = append(connected_players[i].Hand[:i], connected_players[i].Hand[i+1:]...)
			break
		}
	}
	played_cards = append(played_cards, *cardPlayed)
	checkIfAllPlayed()
}

func checkIfAllPlayed() {
	if len(played_cards) == len(connected_players) {
		message := &Message{
			Type:    "AllPlayed",
			Content: "All players have played",
		}
		broadcastMessage(*message)
	}
}

func writeAllPlayedCards() {
	playedCards := &PlayedCardsJSON{played_cards}

	message := &Message{
		Type:    "SendAllPlayedCards",
		Content: playedCards,
	}

	broadcastMessage(*message)
}
