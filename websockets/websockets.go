package websockets

import (
	"net/http"
	"log"
	"encoding/base64"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Lobby struct {
	LobbyName	string
	ClientCount	int64
	Clients		map[string]*Client
}

type Client struct {
	Conn		websocket.Conn
	UserName	string			`json:"user_name"`
}

type WSEvent struct {
	Event		string  	`json:"event"`
	LobbyName	string		`json:"lobby_name"`
	Data		[]string	`json:"data"`
}

var lobbyPool = make(map[string]*Lobby)

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

	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var websocketEvent WSEvent
		json.Unmarshal(msg, &websocketEvent)

		var payload []byte
		switch websocketEvent.Event {
		case "userConnect":
			addToClientPool(websocketEvent, conn)
			payload = getConnectedUserNames(websocketEvent)
		case "userDisconnect":
			removeFromClientPool(websocketEvent.LobbyName, websocketEvent.Data[0])
			payload = getConnectedUserNames(websocketEvent)
		}

		log.Println(lobbyPool[websocketEvent.LobbyName])
		dispatchMsgToPool(websocketEvent, payload, msgType)
	}
}

func addToClientPool(event WSEvent, conn *websocket.Conn) {
	if lobbyPool[event.LobbyName] == nil {
		createNewLobbyPool(event.LobbyName)
	}

	workingLobby := lobbyPool[event.LobbyName]

	client := &Client {
		Conn:		*conn,
		UserName:	event.Data[0],
	}
	workingLobby.Clients[client.UserName] = client
	workingLobby.ClientCount += 1
}

func createNewLobbyPool(lobbyHash string) {
	lobbyName := make([]byte, base64.StdEncoding.DecodedLen(len(lobbyHash)))
	_, err := base64.StdEncoding.Decode(lobbyName, []byte(lobbyHash))
	if err != nil {
		return
	}

	lobbyPool[lobbyHash] = &Lobby{
		LobbyName:		string(lobbyName),
		ClientCount: 	0,
		Clients:		make(map[string]*Client),
	}
}

func dispatchMsgToPool(websocketEvent WSEvent, payload []byte, msgType int) {
	workingLobby := lobbyPool[websocketEvent.LobbyName]

	for _, client := range workingLobby.Clients {
		if err := client.Conn.WriteMessage(msgType, payload); err != nil {
			log.Println(err)
			return
		}
	}
}