package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WsPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type WsMessage struct {
	Payload *WsPayload
	Conn    *WebSocketConnection
}

type WsJsonResponse struct {
	Action        string   `json:"action"`
	Message       string   `json:"message"`
	ConnectedUser []string `json:"connected_users"`
	SessionID     string   `json:"session_id"`
}

type WebSocketConnection struct {
	*websocket.Conn
	SessionID string
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
