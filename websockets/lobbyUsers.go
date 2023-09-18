package websockets

import (
	"encoding/json"
	"log"
)

func (event *WSEvent) getConnectedUserNames() ([]byte, error) {
	workingLobby := lobbyPool[event.LobbyHash]
	responseEvent := &WSEvent{}
	responseEvent.Event = event.Event
	for _, client := range workingLobby.Clients {
		responseEvent.Data = append(responseEvent.Data, client.UserName)
	}

	payload, err := json.Marshal(responseEvent)
	if err != nil {
		log.Print(err)
	}

	return payload, err
}