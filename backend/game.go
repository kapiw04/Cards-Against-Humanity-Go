package main

import (
	"strconv"
	"sync"
)

type GameManager struct {
	Number_of_players int
	Players           map[string]*Player
}

type Player struct {
	ID string
	// Conn *websocket.Conn
	Hand []Card
}

var game_manager_instance *GameManager
var once sync.Once

func New() *GameManager {
	once.Do(func() {
		game_manager_instance = &GameManager{
			Players: make(map[string]*Player),
		}
	})

	return game_manager_instance
}

func cleanGameManager() {
	game_manager := New()
	game_manager.Number_of_players = 0
	game_manager.Players = make(map[string]*Player)
}

func prepareGameManager(number_of_players int) {
	cleanGameManager()
	game_manager := New()
	game_manager.Number_of_players = number_of_players
	game_manager.Players = make(map[string]*Player, number_of_players)
}

func startGame() {
	for i := 0; i < game_manager_instance.Number_of_players; i++ {
		player := &Player{
			ID:   "player" + strconv.Itoa(i),
			Hand: []Card{},
		}
		player.init()
		game_manager_instance.Players[player.ID] = player
	}
}

func (player *Player) init() {
	player.Hand = getRandomWhiteCards(5)
}

func runGameLoop(num_players int) {
	prepareGameManager(num_players)
	startGame()
}
