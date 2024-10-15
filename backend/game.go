package main

import (
	"sync"
)

type GameManager struct {
	number_of_players int
	players []Player
}

type Player struct {
	hand []Card
}

var instance *GameManager
var once sync.Once

func New() *GameManager {
	once.Do(func() {
		instance = &GameManager{}
	})

	return instance
}

func prepareGameManager() {
	game_manager := New()
	game_manager.number_of_players = 2
}

func startGame() {
	prepareGameManager()
	game_manager := instance
	for i := 0; i < game_manager.number_of_players; i++ {
		player := Player{}
		player.init()
		game_manager.players = append(game_manager.players, player)
	}
}

func (player *Player) init() {
	player.hand = getRandomWhiteCards(5)
}

func startGameLoop () {
	startGame()
}