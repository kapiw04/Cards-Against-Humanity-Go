package main

import (
	"fmt"
	"sync"
)

type GameManager struct {
	number_of_players int
	players           []Player
}

type Player struct {
	name string
	hand []Card
}

func (player *Player) String() string {
	return fmt.Sprintf("%s: %v", player.name, player.hand)
}

var game_manager_instance *GameManager
var once sync.Once

func New() *GameManager {
	once.Do(func() {
		game_manager_instance = &GameManager{}
	})

	return game_manager_instance
}

func cleanGameManager() {
	game_manager := New()
	game_manager.number_of_players = 0
	game_manager.players = []Player{}
}

func prepareGameManager(number_of_players int) {
	cleanGameManager()
	game_manager := New()
	game_manager.number_of_players = number_of_players
}

func startGame() {
	game_manager := game_manager_instance
	for i := 0; i < game_manager.number_of_players; i++ {
		player := Player{}
		player.init()
		game_manager.players = append(game_manager.players, player)
	}
}

func (player *Player) init() {
	player.hand = getRandomWhiteCards(5)
}

func runGameLoop (num_players int) {
	prepareGameManager(num_players)
	startGame()
}