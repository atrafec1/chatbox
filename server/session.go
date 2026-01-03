// This file will:
// Create a new session (Session ID, net.Conn, *User)

package server

import (
	"net"
	"sync"
	"time"
)

type Session struct {
	id         string
	CreatedAt  time.Time
	LastActive time.Time
	User       *User
	mu         sync.RWMutex
}

func NewSession(user *User, conn net.Conn) *Session {
	return &Session{
		User:       user,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}
}
