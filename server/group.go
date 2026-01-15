package server

import (
	"fmt"
	"strings"
	"sync"
)

type Group struct {
	id       uint
	name     string
	messages chan *Message
	members  map[string]*Session

	mu sync.RWMutex
}

func NewGroup(id uint, name string) *Group {
	g := &Group{
		id:       id,
		name:     name,
		members:  make(map[string]*Session),
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
		g.BroadcastMsg(msg)
	}
}

func (g *Group) ListMembers(s *Session) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var membersList []string

	for _, member := range g.members {
		membersList = append(membersList, member.User.Name)
	}
	if err := s.SendMsg(strings.Join(membersList, ", ")); err != nil {
		return err
	}
	return nil

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

func (g *Group) BroadcastMsg(msg *Message) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, session := range g.members {
		fmt.Printf("Sending to member %v\n", session.User.Name)
		if session.User.id == msg.UserID {
			continue
		}
		go session.SendMsg(fmt.Sprintf("%s: %s", msg.Username, msg.Content))
	}
}

func (g *Group) Enqueue(msg *Message) {
	g.messages <- msg
}
