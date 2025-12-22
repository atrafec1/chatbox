package client

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func (c *Client) ReadMessage() (string, error) {
	message, err := c.Reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading message: %v", err)
	}
	return message, nil
}

func (c *Client) SendMessage(msg string) error {
	if _, err := c.Writer.WriteString(msg); err != nil {
		return fmt.Errorf("error writing message %v ", err)
	}
	return nil
}
