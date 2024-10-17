package main

import (
	"fmt"
	"math/rand/v2"

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

func (card *Card) String() string {
	return fmt.Sprintf("Card: %s", card.Text)
}

func populateCards() {
	db := openDB()
	db.Where("color = ?", "white").Find(&white_cards)
	db.Where("color = ?", "black").Find(&black_cards)
}

func openDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("cah/cah_cards.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
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

	return white_cards[:k]
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
