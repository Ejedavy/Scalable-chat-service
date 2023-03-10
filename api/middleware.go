package api

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"net/http"
)

func (s *Server) WSMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OTP := r.URL.Query().Get("OTP")
		if len(OTP) == 0 {
			http.Error(w, "OTP not provided as query parameter", 403)
			return
		}
		if len(OTP) != 7 {
			http.Error(w, "Wrong OTP", 403)
			return
		}
		username, err := s.redisClient.Get(OTP).Result()
		if !errors.Is(err, redis.Nil) && err != nil {
			http.Error(w, fmt.Sprintf("Error with the redis %s", err), 500)
			return
		}
		if errors.Is(err, redis.Nil) {
			http.Error(w, fmt.Sprintf("Wrong OTP, Unauthorized"), 403)
			return
		}
		r.Header.Set("username", username)
		s.redisClient.Del(OTP)
		handler.ServeHTTP(w, r)
	}
}
