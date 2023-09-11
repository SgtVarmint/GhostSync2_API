package websockets

func removeFromClientPool(lobbyHash string, userName string) {
	lobbyPool[lobbyHash].ClientCount -= 1
	delete(lobbyPool[lobbyHash].Clients, userName)
}