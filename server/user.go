package server

type User struct {
	id   string
	Name string
}

func NewUser(name ...string) User {
	var newUser User
	if len(name) > 0 && name[0] != "" {
		newUser.Name = name[0]
	}
	newUser.Name = getRandomName()
	return newUser
}

func getRandomName() string {
	return "random"
}
