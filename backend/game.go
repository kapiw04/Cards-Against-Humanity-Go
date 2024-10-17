package main

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID   string
	Conn *websocket.Conn
	Hand []Card
}

func startGame() {
	fmt.Println("Started game with " + strconv.Itoa(len(connected_players)) + " players.")
	for i, player := range connected_players {
		fmt.Println("Player " + strconv.Itoa(i) + " has ID " + player.ID)
		connected_players[i].Hand = getRandomWhiteCards(5)
	}
}

func runGameLoop() {
	startGame()
}
