package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ejedavy/scalable-chat-server/utils"
	"github.com/go-redis/redis/v7"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login Called")
	var data AuthBody
	json.NewDecoder(r.Body).Decode(&data)
	username := data.Username
	password := data.Password

	gottenPassword, err := s.redisClient.Get(username).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		http.Error(w, fmt.Sprintf("Error with the redis %s", err), 500)
		return
	}
	if gottenPassword != password {
		http.Error(w, "Wrong username or password", 403)
		return
	}

	OTP := utils.RandomInteger()
	stringOTP := strconv.Itoa(OTP)
	s.redisClient.Set(stringOTP, username, 5*time.Minute)
	json.NewEncoder(w).Encode(AuthResponse{
		Username: username,
		OTP:      OTP,
	})
}
