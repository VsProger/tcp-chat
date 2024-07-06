package core

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Conn     net.Conn
	Name     string
	ChatRoom *ChatRoom
}

func NewClient(conn net.Conn) *Client {
	return &Client{Conn: conn}
}

func (c *Client) Greet() {
	fmt.Fprint(c.Conn, "Welcome to the TCP Chat Server! Please enter your name: ")
	name, _ := bufio.NewReader(c.Conn).ReadString('\n')
	name = strings.TrimSpace(name)
	c.Name = name
	fmt.Fprintf(c.Conn, "Hello %s, you can join a chat room using /join [room_name] or create one with /create [room_name].\n", c.Name)
}

func (c *Client) ReadMessage() (string, error) {
	message, err := bufio.NewReader(c.Conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(message), nil
}

type MessageHandler interface {
	HandleMessage(client *Client, message string)
}
