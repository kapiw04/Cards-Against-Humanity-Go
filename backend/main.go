package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type StartGameRequest struct {
	PlayersNumber int `json:"players_number"`
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	var start_game_request StartGameRequest
	err := json.NewDecoder(r.Body).Decode(&start_game_request)

	if err != nil {
		panic(err)
	}

	runGameLoop(game_manager_instance.Number_of_players)
	w.WriteHeader(http.StatusOK)
}

func kRandomCardsHandler(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("color")
	origin := r.Header.Get("Origin")
	referer := r.Header.Get("Referer")
	userAgent := r.Header.Get("User-Agent")

	fmt.Printf("Request Details:\n- Origin: %s\n- Referer: %s\n- User-Agent: %s\n", origin, referer, userAgent)

	k, err := strconv.Atoi(r.URL.Query().Get("k"))

	if color == "" || err != nil {
		w.Write([]byte("Color should be provided and k should be a number"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cards []Card

	if color == "white" {
		cards = getRandomWhiteCards(k)
	} else if color == "black" {
		cards = getRandomBlackCards(k)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func main() {
	populateCards()

	fmt.Println("Staring server at port :8080")

	// mux := http.NewServeMux()
	http.HandleFunc("/api", kRandomCardsHandler)
	http.HandleFunc("/start", startGameHandler)
	http.HandleFunc("/ws", websocketHandler)
	err := http.ListenAndServe("0.0.0.0:8080", corsMiddleware(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
