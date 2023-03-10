package chat

import (
	"fmt"
	"github.com/ejedavy/scalable-chat-server/utils"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"sync"
)

type connectedClients map[string]*Client
type Manager struct {
	sync.RWMutex
	connectedClients connectedClients
	redisClient      *redis.Client
	config           utils.Config
}

func NewManager(redisClient *redis.Client, config utils.Config) *Manager {
	return &Manager{
		connectedClients: make(map[string]*Client),
		redisClient:      redisClient,
		config:           config,
	}
}

func (m *Manager) OnConnectionHandler(username string, conn *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	client, err := NewClient(username, conn, m.redisClient, m.config)
	if err != nil {
		HandleWebsocketError(conn, err)
	}
	m.connectedClients[username] = client
	err = client.connect()
	if err != nil {
		HandleWebsocketError(conn, err)
	}
	return
}

func (m *Manager) OnDisconnectHandler(username string, conn *websocket.Conn) chan struct{} {
	ch := make(chan struct{})
	client := m.connectedClients[username]
	conn.SetCloseHandler(func(code int, text string) error {
		err := client.Disconnect()
		if err != nil {
			return err
		}
		delete(m.connectedClients, username)
		close(ch)
		return nil
	})

	return ch
}

func (m *Manager) OnMessageFromRedisChannel(username string, conn *websocket.Conn) {
	client := m.connectedClients[username]
	go func() {
		for msg := range client.messageChan {
			channel := msg.Channel
			content := msg.Payload
			conn.WriteJSON(Message{
				Content: content,
				Channel: channel,
			})
		}
	}()
}

func (m *Manager) OnSendUserMessage(username string, conn *websocket.Conn) {
	client := m.connectedClients[username]
	var msg Message
	err := conn.ReadJSON(&msg)
	if err != nil {
		HandleWebsocketError(conn, err)
	}

	if msg.Command != "" {
		switch msg.Command {
		case "SendMessage":
			err = client.redisClient.Publish(msg.Channel, msg.Content).Err()
			if err != nil {
				HandleWebsocketError(conn, err)
				return
			}
		case "Subscribe":
			if msg.Channel == "" {
				err = fmt.Errorf("no channel provided")
				HandleWebsocketError(conn, err)
				return
			}
			err := client.Subscribe(msg.Channel)
			if err != nil {
				HandleWebsocketError(conn, err)
			}
			return
		case "Unsubscribe":
			if msg.Channel == "" {
				err = fmt.Errorf("no channel provided")
				HandleWebsocketError(conn, err)
				return
			}
			err := client.Unsubscribe(msg.Channel)
			if err != nil {
				HandleWebsocketError(conn, err)
			}
			return
		}
	} else {
		err = fmt.Errorf("no command provided")
		HandleWebsocketError(conn, err)
	}
}
