package main

import (
	"net/http"
	"fmt"
	"log"

	"github.com/SgtVarmint/GhostSync2/config"
	"github.com/SgtVarmint/GhostSync2/authentication"
	"github.com/SgtVarmint/GhostSync2/websockets"
)

func main() {
	port := ":5050"
	config, err := config.ParseConfig()
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/v1/auth", func(w http.ResponseWriter, r *http.Request) {
		authentication.Authenticate(w, r, *config)
	})

	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.Socket(w, r)
	})

	server := http.Server {
		Addr: port,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("error running server: %s\n", err)
	}
}