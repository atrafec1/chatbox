package server

import (
	"errors"
	"fmt"

	"chatbox/database"
	"chatbox/domain"
)

func (s *ChatServer) authenticateUser(c *Client) (*User, error) {
	var user *User
	username, err := s.promptUsername(c)
	if err != nil {
		return nil, err
	}

	userExists, err := database.UsernameExists(s.DB, username)
	if err != nil {
		return nil, err
	}

	if userExists {
		user, err = s.loginFlow(c, username)
		if err != nil {
			return nil, fmt.Errorf("login flow failed: %w", err)
		}
	} else {
		user, err = s.registerFlow(c, username)
		if err != nil {
			return nil, fmt.Errorf("registration flow failed: %w", err)
		}
	}
	return user, nil
}

func (s *ChatServer) promptUsername(c *Client) (string, error) {

	if err := c.SendMessage("Enter username: "); err != nil {
		return "", err
	}
	username, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *ChatServer) promptPassword(c *Client) (string, error) {
	if err := c.SendMessage("Password: "); err != nil {
		return "", err
	}
	password, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	return password, nil
}

func (s *ChatServer) registerUser(username, password string) (*User, error) {
	user, err := database.RegisterUser(s.DB, username, password)
	if err != nil {
		return nil, err
	}
	return &User{
		id:      user.ID,
		Name:    user.Username,
		GroupID: user.Group.ID,
	}, nil
}

func (s *ChatServer) loginUser(username, password string) (*User, error) {
	user, err := database.Login(s.DB, username, password)
	if err != nil {
		return nil, err
	}
	return &User{
		id:      user.ID,
		Name:    user.Username,
		GroupID: user.Group.ID,
	}, nil
}

func (s *ChatServer) loginFlow(c *Client, username string) (*User, error) {
	if err := c.SendMessage(fmt.Sprintf("Welcome back %v!", username)); err != nil {
		return nil, err
	}
	for {
		password, err := s.promptPassword(c)
		if err != nil {
			return nil, err
		}

		user, err := s.loginUser(username, password)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidPassword) {
				if err := c.SendMessage("Invalid password. Please try again."); err != nil {
					fmt.Println("failed to send message, closing session:", err)
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("login failed: %w", err)
		}
		return user, nil
	}
}

func (s *ChatServer) registerFlow(c *Client, username string) (*User, error) {
	if err := c.SendMessage(fmt.Sprintf("Welcome to chatbox %v!", username)); err != nil {
		return nil, err
	}
	password, err := s.promptPassword(c)
	if err != nil {
		return nil, err
	}
	user, err := s.registerUser(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return user, nil
}
