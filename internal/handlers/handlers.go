package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
	"stevematdavies/websockets/chat/helpers"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type ConnectedUser struct {
	Username  string `json:"username"`
	UserColor string `json:"userColor"`
}

//WsJSONResponse defines the response sent back from the websocket
type WsJSONResponse struct {
	Action         string          `json:"action"`
	Message        string          `json:"message"`
	MessageType    string          `json:"messageType"`
	MessageColor   string          `json:"messageColor"`
	ConnectedUsers []ConnectedUser `json:"connectedUsers"`
}

type WsJSONPayload struct {
	Action   string              `json:"action"`
	Message  string              `json:"message"`
	Username string              `json:"username"`
	Conn     WebSocketConnection `json:"-"`
}

var WsChan = make(chan WsJSONPayload)
var connectedClients = make(map[WebSocketConnection]ConnectedUser)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = view.Execute(w, data, nil); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func Home(w http.ResponseWriter, _ *http.Request) {
	if err := renderPage(w, "home.jet", nil); err != nil {
		log.Println(err)
	}
}

// WsEndpoint upgrades the connection to a websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to endpoint")

	conn := WebSocketConnection{Conn: ws}
	connectedClients[conn] = ConnectedUser{
		UserColor: helpers.GetRandomChatColor(),
	}

	if err = ws.WriteJSON(WsJSONResponse{
		Message: `<em><small>Connected to server</small></em>`,
	}); err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection) {

	var payload WsJSONPayload

	defer func() {
		if r := recover(); r != nil {
			log.Println("Error: ", fmt.Sprintf("%v", r))
		}
	}()

	for {
		if err := conn.ReadJSON(&payload); err != nil {
			log.Println("Error: ", fmt.Sprintf("%v", err))
		} else {
			payload.Conn = *conn
			WsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var wjr WsJSONResponse
	for {
		e := <-WsChan
		fmt.Println(e)
		switch e.Action {

		case "username":
			connectedClients[e.Conn] = ConnectedUser{e.Username, connectedClients[e.Conn].UserColor}
			wjr.Action = "list_users"
			wjr.ConnectedUsers = getConnectedUsers()
			broadCastToAll(wjr)

		case "broadcast":
			wjr.Action = "broadcast"
			wjr.Message = fmt.Sprintf("<span style=\"color:%s\"><strong>%s</strong></span> %s ", connectedClients[e.Conn].UserColor, e.Username, e.Message)
			broadCastToAll(wjr)

		case "left":
			wjr.Action = "list_users"
			delete(connectedClients, e.Conn)
			wjr.ConnectedUsers = getConnectedUsers()
			broadCastToAll(wjr)
		}

	}
}

func getConnectedUsers() []ConnectedUser {
	var clients []ConnectedUser
	for _, c := range connectedClients {
		if c.Username != "" {
			clients = append(clients, c)
		}
	}
	sort.Slice(clients, func(i, j int) bool {
		return clients[i].Username < clients[j].Username
	})
	return clients
}

func broadCastToAll(r WsJSONResponse) {
	for client := range connectedClients {
		if err := client.WriteJSON(r); err != nil {
			log.Println("Websocket Error")
			_ = client.Close()
			delete(connectedClients, client)
		}
	}
}
