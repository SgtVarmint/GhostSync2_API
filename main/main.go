package main

import (
	"net/http"
	"fmt"

	"github.com/SgtVarmint/GhostSync2/authentication"
	"github.com/SgtVarmint/GhostSync2/websockets"
)

func main() {
	port := ":5050"

	http.HandleFunc("/v1/auth", func(w http.ResponseWriter, r *http.Request) {
		authJson := authentication.Authenticate(r)
		authentication.SendAuthResponse(w, authJson)
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