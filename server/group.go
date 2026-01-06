package server

import "sync"

type Group struct {
	id      uint
	name    string
	members map[uint]*Session
	mu      sync.RWMutex
}

func NewGroup(id uint, name string) *Group {
	return &Group{
		id:      id,
		name:    name,
		members: make(map[uint]*User),
	}
}

// unexported helpers — server only
func (g *Group) Add(s *Session) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.members[u.id] = u
}

func (g *Group) Remove(s *User) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.members, u.id)
}

func (g *Group) BroadcastMsg(msg string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, user := range g.members {
		user.SendMessage(msg)
	}
}
