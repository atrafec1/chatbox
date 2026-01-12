package server

import (
	"fmt"
	"sync"
)

type Group struct {
	id       uint
	name     string
	messages chan *Message
	members  map[uint]*Session

	mu sync.RWMutex
}

func NewGroup(id uint, name string) *Group {
	g := &Group{
		id:       id,
		name:     name,
		members:  make(map[uint]*Session),
		messages: make(chan *Message, 128),
	}
	go g.distributeMessages()
	return g
}

func (g *Group) Close() {
	close(g.messages)
}

func (g *Group) distributeMessages() {
	for msg := range g.messages {
		formatted := fmt.Sprintf("%s: %s", msg.Username, msg.Content)
		g.BroadcastMsg(formatted)
	}
}

// unexported helpers — server only
func (g *Group) Add(s *Session) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.members[s.id] = s
}

func (g *Group) Remove(s *Session) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.members, s.id)
}

func (g *Group) BroadcastMsg(msg string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, session := range g.members {
		fmt.Printf("Sending to member %v\n", session.User.Name)
		go session.SendMsg(msg)
	}
}

func (g *Group) Enqueue(msg *Message) {
	g.messages <- msg
}
