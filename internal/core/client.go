package core

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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
	f, err := os.Open("greetMessage.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(lines)-1; i++ {
		fmt.Fprintln(c.Conn, lines[i])
	}

	if len(lines) > 0 {
		fmt.Fprint(c.Conn, lines[len(lines)-1]+" ")
	}

	name, _ := bufio.NewReader(c.Conn).ReadString('\n')
	name = strings.TrimSpace(name)
	c.Name = name
	fmt.Fprintf(c.Conn, "Hello %s. Use the /help command to get a list of commands\n", c.Name)
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
