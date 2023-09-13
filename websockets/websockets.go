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
	LobbyHash	string
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
		websocketEvent.getLobbyHash()

		payload, err := ParseEvent(websocketEvent, conn)
		if err != nil {
			panic(err)
		}

		dispatchMsgToLobby(websocketEvent.LobbyHash, payload, msgType)
	}
}

func (event *WSEvent) getLobbyHash() {
	event.LobbyHash = base64.StdEncoding.EncodeToString([]byte(event.LobbyName))
} 

func ParseEvent(websocketEvent WSEvent, conn *websocket.Conn) ([]byte, error){
	var payload []byte
	var err error

	switch websocketEvent.Event {
	case "userConnect":
		websocketEvent.AddToClientPool(conn)
		payload, err = websocketEvent.getConnectedUserNames()
	case "userDisconnect":
		websocketEvent.removeFromClientPool()
		payload, err = websocketEvent.getConnectedUserNames()
	default:
		payload, err = json.Marshal(websocketEvent)
		if err != nil {
			panic(err)
		}
	}

	return payload, err
}

func (event *WSEvent) CreateNewLobby() (*Lobby){
	return &Lobby{
		LobbyName:		event.LobbyName,
		ClientCount: 	0,
		Clients:		make(map[string]*Client),
	}
}

func dispatchMsgToLobby(lobbyHash string, payload []byte, msgType int) {
	workingLobby := lobbyPool[lobbyHash]

	for _, client := range workingLobby.Clients {
		if err := client.Conn.WriteMessage(msgType, payload); err != nil {
			log.Println(err)
			return
		}
	}
}