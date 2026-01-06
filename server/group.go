package server

import "sync"

type Group struct {
	id      uint
	name    string
	members map[uint]*User
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
func (g *Group) Add(u *User) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.members[u.id] = u
}

func (g *Group) Remove(u *User) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.members, u.id)
}
