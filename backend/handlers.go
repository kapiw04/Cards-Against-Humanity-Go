package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getPlayerFromAddr(addr string) (*Player, int, error) {
	for i, player := range connected_players {
		if player.Addr == addr {
			return &connected_players[i], i, nil // Return pointer to the actual player in the slice
		}
	}
	return nil, -1, fmt.Errorf("Player with address %s not found", addr)
}

// @Summary Start the game
// @Description Start the game
// @Tags game
// @Success 200 {string} string "OK"
// @Router /start [post]
func startGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	runGameLoop()
	w.WriteHeader(http.StatusOK)
}

// @Summary Get Hand of the player
// @Description Get the hand of the player
// @Tags player
// @Produce json
// @Success 200 {object} []Card
// @Failure 403 {string} string "Forbidden"
// @Router /hand [get]
func getHandHandler(w http.ResponseWriter, r *http.Request) {
	player, _, err := getPlayerFromAddr(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player.Hand)
}

// @Summary Get current black card
// @Description Retrieve the current black card
// @Tags game
// @Produce json
// @Success 200 {object} Card
// @Router /blackcard [get]
func getBlackCardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentBlackCard)
}

// @Summary Player plays a card
// @Description Play a card by providing the card index in the header
// @Tags game
// @Param CardIndex header int true "Index of the card to play"
// @Success 202 {string} string "Card accepted"
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Forbidden"
// @Router /card/play [post]
func cardPlayedHandler(w http.ResponseWriter, r *http.Request) {
	card_index, err := strconv.Atoi(r.Header.Get("CardIndex"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	player, player_index, err := getPlayerFromAddr(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(err)
	}
	if card_index < 0 || len(player.Hand) <= card_index {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	card := player.Hand[card_index]

	connected_players[player_index].Hand = append(player.Hand[:card_index],
		player.Hand[card_index+1:]...) // remove card from hand

	played_cards[card] = connected_players[player_index]
	if checkIfAllPlayed() {
		fmt.Println("Ready to move to round 2")
	}

	w.WriteHeader(http.StatusAccepted)
}

// @Summary Get all played cards
// @Description Retrieve all the cards that have been played
// @Tags game
// @Produce json
// @Success 200 {array} Card
// @Router /cards/played [get]
func getAllPlayedCardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	keys := make([]Card, 0, len(played_cards)) // why has len 1
	playerAddresses := make([]string, 0, len(played_cards))

	for k, v := range played_cards {
		keys = append(keys, k)
		playerAddresses = append(playerAddresses, v.Addr)
	}
	json.NewEncoder(w).Encode(keys)
	json.NewEncoder(w).Encode(playerAddresses)
}
