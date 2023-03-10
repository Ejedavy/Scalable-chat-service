package chat

import "github.com/gorilla/websocket"

func HandleWebsocketError(conn *websocket.Conn, err error) {
	msg := Message{
		Error: err,
	}
	conn.WriteJSON(msg)
}

type Message struct {
	Command string `json:"command,omitempty"`
	Content string `json:"content,omitempty"`
	Channel string `json:"channel,omitempty"`
	Error   error  `json:"error,omitempty"`
}
