package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/MrBorisT/is-it-you/server/internal/config"
	"github.com/MrBorisT/is-it-you/server/internal/game"
	"github.com/MrBorisT/is-it-you/server/internal/transport/ws"
)

func main() {
	configPath := flag.String("config", "config.cfg", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("load config:", err)
	}

	g := game.NewGame(cfg.Game)
	handler := ws.NewHandler(g)

	go g.Loop(handler.Broadcast)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ws", handler.HandleWS)

	log.Println("server listening on", cfg.Addr())
	log.Fatal(http.ListenAndServe(cfg.Addr(), mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
