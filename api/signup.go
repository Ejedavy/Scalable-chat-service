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

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignUp Called")
	var data AuthBody
	json.NewDecoder(r.Body).Decode(&data)
	username := data.Username
	password := data.Password

	if len(username) < 5 || len(password) < 5 {
		http.Error(w, "Username and Password should be at least 5 characters", 400)
		return
	}

	gottenPassword, err := s.redisClient.Get(username).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		http.Error(w, fmt.Sprintf("Error with the redis %s", err), 500)
		return
	}
	if len(gottenPassword) != 0 {
		http.Error(w, "Username already exists", 400)
		return
	}

	s.redisClient.Set(username, password, 0)
	OTP := utils.RandomInteger()
	stringOTP := strconv.Itoa(OTP)
	s.redisClient.Set(stringOTP, username, 5*time.Minute)
	json.NewEncoder(w).Encode(AuthResponse{
		Username: username,
		OTP:      OTP,
	})
}
