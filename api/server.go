package api

import (
	"github.com/ejedavy/scalable-chat-server/chat"
	"github.com/ejedavy/scalable-chat-server/utils"
	"github.com/go-redis/redis/v7"
)

type Server struct {
	redisClient *redis.Client
	config      utils.Config
	manager     *chat.Manager
}

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Username string `json:"username"`
	OTP      int    `json:"OTP"`
}

func NewServer(client *redis.Client, config utils.Config) (*Server, error) {
	return &Server{
		redisClient: client,
		config:      config,
		manager:     chat.NewManager(client, config),
	}, nil
}
