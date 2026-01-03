package server

type Group struct {
	id      string
	name    string
	members map[string]*User
}

func newGroup(id, name string) *Group {
	return &Group{
		id:      id,
		name:    name,
		members: make(map[string]*User),
	}
}

// unexported helpers — server only
func (g *Group) add(u *User) {
	g.members[u.id] = u
}

func (g *Group) remove(u *User) {
	delete(g.members, u.id)
}
