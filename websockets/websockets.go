package websockets

import (
	"net/http"
	"log"

	"github.com/gorilla/websocket"
)

var clients []websocket.Conn

var upgrader = websocket.Upgrader {
	ReadBufferSize:		1024,
	WriteBufferSize:	1024,
}

func Socket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func (r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	clients = append(clients, *conn)
	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%s, %s\n", messageType, string(msg))

		for _, client := range clients {
			if err := client.WriteMessage(messageType, msg); err != nil {
				log.Println(err)
				return
			}
		}
	}
}