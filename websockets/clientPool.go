package websockets

import (
	"github.com/gorilla/websocket"
)

func (event *WSEvent) removeFromClientPool() {
	userName := event.Data[0]
	workingLobby := lobbyPool[event.LobbyHash]
	workingLobby.ClientCount -= 1
	delete(workingLobby.Clients, userName)
}

func (event *WSEvent) AddToClientPool(conn *websocket.Conn) {
	if lobbyPool[event.LobbyHash] == nil {
		lobbyPool[event.LobbyHash] = event.CreateNewLobby()
	}

	workingLobby := lobbyPool[event.LobbyHash]

	client := &Client {
		Conn:		*conn,
		UserName:	event.Data[0],
	}
	workingLobby.Clients[client.UserName] = client
	workingLobby.ClientCount += 1
}