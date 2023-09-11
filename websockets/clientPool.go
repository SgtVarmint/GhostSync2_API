package websockets

func removeFromClientPool(lobbyHash string) {
	lobbyPool[lobbyHash].ClientCount -= 1
}