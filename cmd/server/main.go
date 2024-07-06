package main

import (
	"fmt"
	"net"
	"os"
	"tcp-chat/internal/chat"
	"tcp-chat/internal/core"
	"tcp-chat/internal/handlers"
)

func main() {
	args := os.Args
	port := "8989"

	if len(args) > 2 {
		fmt.Println("[USAGE]: go run ./cmd/server/ $port")
		return
	} else if len(args) == 2 {
		port = args[1]
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	chat.InitializeChatRooms()

	handler := &handlers.ChatMessageHandler{}
	
	fmt.Println("Launching server...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := core.NewClient(conn)
		go handleConnection(client, handler)
	}
}

func handleConnection(client *core.Client, handler core.MessageHandler) {
	defer client.Conn.Close()

	client.Greet()

	for {
		message, err := client.ReadMessage()
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		handler.HandleMessage(client, message)
	}
}
