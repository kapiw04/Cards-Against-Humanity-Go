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

func getHandHandler(w http.ResponseWriter, r *http.Request) {
	player, _, err := getPlayerFromAddr(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player.Hand)
}

func getBlackCardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentBlackCard)
}

func CardPlayedHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "Card played successfully")
	w.WriteHeader(http.StatusAccepted)
}

func getAllPlayedCardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	keys := make([]Card, len(played_cards))
	players := make([]Player, len(played_cards))

	for k, v := range played_cards {
		keys = append(keys, k)
		players = append(players, v)
	}
	json.NewEncoder(w).Encode(keys)
	json.NewEncoder(w).Encode(players)
}
