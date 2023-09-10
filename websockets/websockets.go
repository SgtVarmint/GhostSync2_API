package websockets

import (
	"net/http"
	"log"
	"encoding/base64"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type lobby struct {
	LobbyName	string
	ClientCount	int64
	Clients		[]Client
}

type Client struct {
	Conn		websocket.Conn
	UserName	string			`json:"user_name"`
}

type UserConnect struct {
	Event		string  `json:"event"`
	LobbyName	string	`json:"lobby_name"`
	UserName	string	`json:"user_name"`
}

var lobbyPool = make(map[string]*lobby)

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

		var WSEvent UserConnect
		json.Unmarshal(msg, &WSEvent)
	
		switch WSEvent.Event {
		case "userConnect":
			addToClientPool(WSEvent, conn)
		case "userDisconnect":
			removeFromClientPool(WSEvent.LobbyName)
		}

		dispatchMsgToPool(WSEvent, msgType)
	}
}

func addToClientPool(WSEvent UserConnect, conn *websocket.Conn) {
	if lobbyPool[WSEvent.LobbyName] == nil {
		createNewLobbyPool(WSEvent.LobbyName)
	}
	
	workingLobby := lobbyPool[WSEvent.LobbyName]

	client := &Client {
		Conn:		*conn,
		UserName:	WSEvent.UserName,
	}
	workingLobby.Clients = append(workingLobby.Clients, *client)
	workingLobby.ClientCount += 1
}

func createNewLobbyPool(lobbyHash string) {
	lobbyName := make([]byte, base64.StdEncoding.DecodedLen(len(lobbyHash)))
	_, err := base64.StdEncoding.Decode(lobbyName, []byte(lobbyHash))
	if err != nil {
		return
	}

	lobbyPool[lobbyHash] = &lobby{
		LobbyName:		string(lobbyName),
		ClientCount: 	0,
	}
}

func removeFromClientPool(lobbyHash string) {
	lobbyPool[lobbyHash].ClientCount -= 1
}

func dispatchMsgToPool(WSEvent UserConnect, msgType int) {
	workingLobby := lobbyPool[WSEvent.LobbyName]

	payload, err := json.Marshal(WSEvent)
	if err != nil {
		panic(err)
	}

	for _, client := range workingLobby.Clients {
		if err = client.Conn.WriteMessage(msgType, []byte(payload)); err != nil {
			log.Println(err)
			return
		}
	}
}