package api

import (
	"github.com/ejedavy/scalable-chat-server/chat"
	"github.com/gorilla/websocket"
	"net/http"
)

func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		chat.HandleWebsocketError(conn, err)
	}

	manager := s.manager
	manager.OnConnectionHandler(username, conn)
	ch := manager.OnDisconnectHandler(username, conn)
	manager.OnMessageFromRedisChannel(username, conn)

	for {
		select {
		case <-ch:
			return
		default:
			manager.OnSendUserMessage(username, conn)
		}
	}

}
