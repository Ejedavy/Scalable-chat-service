package main

import (
	"fmt"
	"github.com/ejedavy/scalable-chat-server/api"
	"github.com/ejedavy/scalable-chat-server/utils"
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

var (
	redisClient *redis.Client
	server      *api.Server
	config      utils.Config
)

func runSetup() {
	var err error
	config, err = utils.NewConfig()
	if err != nil {
		log.Fatal("Could not start server")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
		DB:   0,
	})
	if redisClient == nil {
		log.Fatalf("Redis is not started on %s", config.RedisAddress)
	}
	redisClient.SAdd(config.GeneralChannels, "eje")
	server, err = api.NewServer(redisClient, config)
	if err != nil {
		log.Fatal("Could not start server")
	}
}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/community", server.WSMiddleware(server.WebSocketHandler))
	mux.HandleFunc("/signup", server.SignUpHandler)
	mux.HandleFunc("/login", server.LoginHandler)
}

func main() {
	runSetup()
	defer redisClient.Close()
	mux := http.NewServeMux()
	registerHandlers(mux)
	fmt.Println("Listening at", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, mux))
}
