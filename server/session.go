// This file will:
// Create a new session (Session ID, net.Conn, *User)

package server

import (
	"sync"
	"time"
)

type Session struct {
	id         string
	Client     *Client
	CreatedAt  time.Time
	LastActive time.Time
	User       *User
	mu         sync.RWMutex
}

func NewSession(user *User, client *Client) *Session {
	return &Session{
		User:       user,
		Client:     client,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}
}

// Server sending to client
func (s *Session) SendMsg(msg string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.Client.SendMessage(msg); err != nil {
		return err
	}
	return nil
}

func (s *Session) ReadMsg() (string, error) {
	return s.Client.ReadMessage()
}

func (s *Session) UpdateLastActive() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastActive = time.Now()
}
