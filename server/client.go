package server

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

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
}

func (c *Client) ReadMessage() (string, error) {
	message, err := c.Reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading message: %v", err)
	}
	return message, nil
}

func (c *Client) SendMessage(msg string) error {
	if _, err := c.Writer.WriteString(msg + "\n"); err != nil { // add newline so ReadLine works client-side
		return fmt.Errorf("error writing message: %v", err)
	}
	if err := c.Writer.Flush(); err != nil {
		return fmt.Errorf("error flushing message: %v", err)
	}
	return nil
}
