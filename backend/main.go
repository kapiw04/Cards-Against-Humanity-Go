package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Color string

type Card struct {
	ID    uint `gorm:"primaryKey"`
	Text  string
	Color string
}

var white_cards []Card
var black_cards []Card

func populateCards() {
	db, err := gorm.Open(sqlite.Open("cah/cah_cards.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Where("color = ?", "white").Find(&white_cards)
	db.Where("color = ?", "black").Find(&black_cards)
}

func getRandomWhiteCards(k int) []Card {
	cards_len := len(white_cards)
	rand.Shuffle(cards_len, func(i, j int) {
		white_cards[i], white_cards[j] = white_cards[j], white_cards[i]
	})

	if k > cards_len {
		k = cards_len
		fmt.Println("WARNING: provided k was greater than cards length")
	}

	return white_cards[:k+1]
}

func getRandomBlackCards(k int) []Card {
	cards_len := len(black_cards)
	rand.Shuffle(cards_len, func(i, j int) {
		black_cards[i], black_cards[j] = black_cards[j], black_cards[i]
	})

	if k > cards_len {
		k = cards_len
		fmt.Println("WARNING: provided k was greater than cards length")
	}

	return black_cards[:k]
}

func kRandomCardsHandler(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("color")
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

	http.HandleFunc("/api/", kRandomCardsHandler)
	http.ListenAndServe(":8080", nil)
}