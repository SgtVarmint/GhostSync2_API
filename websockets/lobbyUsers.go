package websockets

import (
	"encoding/json"
	"log"
)

type ConnectedUsers struct {
	Event		string
	Users		[]string
}

func getConnectedUserNames(event WSEvent) ([]byte){
	workingLobby := lobbyPool[event.LobbyName]
	connUsers := &ConnectedUsers{}
	connUsers.Event = event.Event
	for _, client := range workingLobby.Clients {
		connUsers.Users = append(connUsers.Users, client.UserName)
	}

	payload, err := json.Marshal(connUsers)
	if err != nil {
		log.Print(err)
	}

	return payload
}