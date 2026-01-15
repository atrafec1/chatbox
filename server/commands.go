package server

import (
	. "chatbox/domain"
	"fmt"
	"strings"
)

type CommandHandler func(args []string, session *Session) error

type Command struct {
	Name        string
	Handler     CommandHandler
	Usage       string
	Description string
	MinArgs     int
}

var Commands = map[string]Command{}

func init() {
	Commands["hello"] = hello
	Commands["groups"] = listGroupings
	Commands["commands"] = commands
	Commands["listmembers"] = CMDListGroupMembers
	// Commands["join group"] = joinGroup
}

func IsCommand(msg string) bool {
	if strings.HasPrefix(msg, "/") {
		return true
	}
	return false
}

var commands = Command{
	Name:        "commands",
	Handler:     listCommands,
	Usage:       "/commands",
	Description: "List all available commands for your authentication level",
	MinArgs:     0,
}

func ParseCommand(command string) (*Command, []string, error) {
	// /join <group_name>
	if !strings.HasPrefix(command, "/") {
		return nil, nil, ErrNotCommand
	}
	cmd_without_prefix := command[1:]
	parts := strings.Split(cmd_without_prefix, " ")

	cmdName := parts[0]
	arguments := parts[1:]

	cmd, exists := Commands[cmdName]
	if !exists {
		return nil, nil, ErrNotCommand
	}

	if cmd.MinArgs > len(arguments) {
		return nil, nil, ErrInvalidCommandArgs
	}
	fmt.Printf("%v: arguments", arguments)
	return &cmd, arguments, nil
}

var listGroupings = Command{
	Name:        "groups",
	Handler:     listGroups,
	Usage:       "/groups",
	Description: "Get a list of all the groups",
	MinArgs:     0,
}

func listCommands(args []string, s *Session) error {
	var commandsList []string
	for _, cmd := range Commands {
		commandsList = append(commandsList, fmt.Sprintf("%s: %s", cmd.Usage, cmd.Description))
	}
	if err := s.SendMsg("Available commands:\n" + strings.Join(commandsList, "\n")); err != nil {
		return err
	}
	return nil
}
func joinGroup(args []string, s *Session) error {
	//Were gonna use the argument (should only be one groupname) #add user to group on the server and database level
	return nil
}

func listGroups(args []string, s *Session) error {
	var groups []string
	for i := range 5 {
		fake_group := fmt.Sprintf("Group %d", i+1)
		groups = append(groups, fake_group)
	}
	if err := s.SendMsg("Available groups:\n" + strings.Join(groups, "\n")); err != nil {
		return err
	}
	return nil
}

var hello = Command{
	Name:        "hello",
	Handler:     HelloWorld,
	Usage:       "/hello",
	Description: "Say hello to your group",
	MinArgs:     0,
}

func HelloWorld(args []string, s *Session) error {
	msg := &Message{
		UserID:   s.User.id,
		Username: s.User.Name,
		Content:  "Hello from command",
		GroupID:  s.User.GroupID,
	}
	if err := s.Server.routeMessage(msg); err != nil {
		return err
	}
	return nil
}

func (s *Session) ExecuteCommand(args []string, cmd *Command) error {
	if err := cmd.Handler(args, s); err != nil {
		return err
	}
	return nil
}

var CMDListGroupMembers = Command{
	Name:        "list group members",
	Handler:     listGroupMembers,
	Description: "List group members given a specified group name",
	Usage:       "/listmembers <groupname>",
	MinArgs:     1,
}

func listGroupMembers(args []string, s *Session) error {
	groupName := args[0]
	group, err := s.Server.getGroupByName(groupName)
	if err != nil {
		return err
	}
	if err := group.ListMembers(s); err != nil {
		return err
	}
	return nil
}
