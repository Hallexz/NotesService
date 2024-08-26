package auntification

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

func NewSession(redisClient *redis.Client) *Session {
	return &Session{redisClient: redisClient}
}

type Session struct {
	redisClient *redis.Client
}

func (s *Session) CreateSession(user *User) (string, error) {
	sessionID := generateSessionID()
	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	err = s.redisClient.HMSet(sessionID, map[string]interface{}{"user": userJSON}).Err()
	if err != nil {
		return "", err
	}
	s.redisClient.Expire(sessionID, time.Hour*24)
	return sessionID, nil
}

func (s *Session) GetUser(sessionID string) (*User, error) {
	user, err := s.redisClient.HGet(sessionID, "user").Result()
	if err != nil {
		return nil, err
	}
	var u User
	err = json.Unmarshal([]byte(user), &u)
	if err != nil {
		return nil, err
	}
	return &u, err
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
