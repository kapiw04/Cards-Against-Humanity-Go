package main

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

type Player struct {
	Addr string
	Conn *websocket.Conn
	Hand []Card
}

var currentBlackCard Card

func startGame() {
	fmt.Println("Started game with " + strconv.Itoa(len(connected_players)) + " players.")
	for i := range connected_players {
		connected_players[i].Hand = getRandomWhiteCards(5)
	}

	played_cards = make(map[Card]Player, len(connected_players))
	currentBlackCard = getRandomBlackCards(1)[0]
}

func runGameLoop() {
	startGame()
}
