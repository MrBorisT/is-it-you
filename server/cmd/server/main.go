package main

import (
	"log"
	"net/http"

	"github.com/MrBorisT/is-it-you-server/internal/game"
	"github.com/MrBorisT/is-it-you-server/internal/transport/ws"
)

func main() {
	g := game.NewGame()
	handler := ws.NewHandler(g)

	go g.Loop(handler.Broadcast)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ws", handler.HandleWS)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
