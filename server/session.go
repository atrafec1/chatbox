// This file will:
// Create a new session (Session ID, net.Conn, *User)

package server

import (
	"net"
)

type Session struct {
	Connection net.Conn
	User       User
}

func NewSession(user User, conn net.Conn) *Session {
	return &Session{
		Connection: conn,
		User:       user,
	}
	var session Session
	session.Connection = conn
	session.User = user
	return &session
}
