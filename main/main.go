package main

import (
	"net/http"
	"fmt"

	"github.com/SgtVarmint/GhostSync2/config"
	"github.com/SgtVarmint/GhostSync2/authentication"
)

func main() {
	port := ":5050"
	config, err := config.ParseConfig()
	
	if err != nil {
		fmt.Errorf("Uh oh! :(")
	}

	http.HandleFunc("/v1/auth", func(w http.ResponseWriter, r *http.Request) {
		authentication.Authenticate(w, r, *config)
	})

	server := http.Server {
		Addr: port,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("error running server: %s\n", err)
	}

	fmt.Printf("Starting server on port: %s...\n", port)
}