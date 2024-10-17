package main

import (
	"fmt"
	"net/http"
)

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	runGameLoop()
	w.WriteHeader(http.StatusOK)
}

func main() {
	populateCards()

	fmt.Println("Staring server at http://localhost:8080")

	http.HandleFunc("/start", startGameHandler)
	http.HandleFunc("/ws", websocketHandler)
	err := http.ListenAndServe("0.0.0.0:8080", corsMiddleware(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // temporary solution
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
