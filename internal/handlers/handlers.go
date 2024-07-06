package handlers

import (
	"fmt"
	"strings"
	"tcp-chat/internal/chat"
	"tcp-chat/internal/core"
	"tcp-chat/internal/encryption"
	"time"
)

type ChatMessageHandler struct{}

var encryptionKey = []byte("aes256key-32characterslongpasswo")

func (h *ChatMessageHandler) HandleMessage(client *core.Client, message string) {
	if strings.HasPrefix(message, "/") {
		h.handleCommand(client, message)
		return
	}

	h.processChatMessage(client, message)
}

func (h *ChatMessageHandler) handleCommand(client *core.Client, command string) {
	command = strings.TrimSpace(command)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Fprintf(client.Conn, "Unknown command. Type /help for command list.\n")
		return
	}

	switch parts[0] {
	case "/help":
		h.help(client)
	case "/create":
		if len(parts) < 2 {
			h.help(client)
			return
		}
		chatName := strings.Join(parts[1:], " ")
		h.createChatRoom(client, chatName)
	case "/join":
		if len(parts) < 2 {
			h.help(client)
			return
		}
		chatName := strings.Join(parts[1:], " ")
		h.joinChatRoom(client, chatName)
	case "/users":
		if client.ChatRoom == nil {
			fmt.Fprintf(client.Conn, "You are not in a chat room.\n")
		} else {
			h.listUsers(client)
		}
	case "/logout":
		if client.ChatRoom == nil {
			fmt.Fprintf(client.Conn, "You are not in a chat room.\n")
		} else {
			h.leaveChatRoom(client)
		}
	// case "/kick":
	// 	if len(parts) < 2 {
	// 		fmt.Fprintf(client.Conn, "Usage: /kick [username]\n")
	// 		return
	// 	}
	// 	username := strings.Join(parts[1:], " ")
	// 	h.kick(client, username)
	// case "/ban":
	// 	if len(parts) < 2 {
	// 		fmt.Fprintf(client.Conn, "Usage: /ban [username]\n")
	// 		return
	// 	}
	// 	username := strings.Join(parts[1:], " ")
	// 	h.ban(client, username)
	default:
		fmt.Fprintf(client.Conn, "Unknown command. Type /help for command list.\n")
	}
}

func (h *ChatMessageHandler) processChatMessage(client *core.Client, message string) {
	if client.ChatRoom == nil {
		fmt.Fprintf(client.Conn, "You are not currently in any chat room. Please join one to start chatting.\n")
		return
	}

	encryptedMessage, err := encryption.Encrypt(encryptionKey, message)
	if err != nil {
		fmt.Fprintf(client.Conn, "Error encrypting message: %s\n", err)
		fmt.Println("Error encrypting message:", err)
		return
	}

	fmt.Println("Encrypted message:", encryptedMessage)
	h.broadcast(client.ChatRoom, client.Name, encryptedMessage)
}

func (h *ChatMessageHandler) broadcast(chatRoom *core.ChatRoom, username, encryptedMessage string) {
	chatRoom.Lock.Lock()
	defer chatRoom.Lock.Unlock()

	currentTime := time.Now().Format("15:04")
	formattedMessage := fmt.Sprintf("%s: [%s] %s\n", currentTime, username, encryptedMessage)
	fmt.Println("Broadcasting encrypted message:", formattedMessage)

	for _, client := range chatRoom.Clients {
		decryptedMessage, err := encryption.Decrypt(encryptionKey, encryptedMessage)
		if err != nil {
			continue
		}

		fmt.Println("Decrypted message:", decryptedMessage)
		finalMessage := fmt.Sprintf("%s: [%s] %s\n", currentTime, username, decryptedMessage)
		if _, err := client.Conn.Write([]byte(finalMessage)); err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}
}

func (h *ChatMessageHandler) help(client *core.Client) {
	helpText := `
Available commands:
/help - Shows help information.
/create [room_name] - Creates a new chat room.
/join [room_name] - Joins an existing chat room.
/users - Shows a list of users(only available in chat)
`
	fmt.Fprintf(client.Conn, helpText)
}

// func (h *ChatMessageHandler) logout(client *core.Client) {
// 	client.ChatRoom = nil
// 	fmt.Fprintf(client.Conn, "You logout from the room.\n")
// 	fmt.Fprintf(client.Conn, "%s has left our chat...\n", username)

// 	normalizedUsername := strings.ToLower(client.Name)

// 	found := false

// 	client.ChatRoom.Lock.Lock()
// 	newClients := []*core.Client{}
// 	for _, c := range client.ChatRoom.Clients {
// 		if strings.ToLower(c.Name) == normalizedUsername {
// 			fmt.Fprintf(c.Conn, "You have been kicked from the room.\n")
// 			c.ChatRoom = nil
// 			client.ChatRoom.KickedUsers[normalizedUsername] = true
// 			found = true
// 		} else {
// 			newClients = append(newClients, c)
// 		}
// 	}
// 	client.ChatRoom.Clients = newClients
// 	client.ChatRoom.Lock.Unlock()

// 	if found {
// 		h.broadcast(client.ChatRoom, "Server", fmt.Sprintf("%s has been kicked from the room.", username))
// 		fmt.Fprintf(client.Conn, "%s has been kicked from the room.\n", username)
// 	}
// }

func (h *ChatMessageHandler) createChatRoom(client *core.Client, name string) {
	chat.ChatRoomsLock.Lock()
	defer chat.ChatRoomsLock.Unlock()

	if _, exists := chat.ChatRooms[name]; exists {
		fmt.Fprintf(client.Conn, "Chat room '%s' already exists.\n", name)
		return
	}

	chat.ChatRooms[name] = core.NewChatRoom(name, client)
	fmt.Fprintf(client.Conn, "Chat room '%s' created. You can join now using '/join %s'.\n", name, name)
}

func (h *ChatMessageHandler) joinChatRoom(client *core.Client, name string) {
	chat.ChatRoomsLock.Lock()
	chatRoom, exists := chat.ChatRooms[name]
	if !exists {
		fmt.Fprintf(client.Conn, "Chat room '%s' does not exist. Create it using '/create %s'.\n", name, name)
		chat.ChatRoomsLock.Unlock()
		return
	}

	if chatRoom.KickedUsers[strings.ToLower(client.Name)] {
		fmt.Fprintf(client.Conn, "You have been kicked from this chat room and cannot join.\n")
		chat.ChatRoomsLock.Unlock()
		return
	}

	chat.ChatRoomsLock.Unlock()

	if client.ChatRoom != nil {
		h.leaveChatRoom(client)
	}

	chatRoom.Lock.Lock()
	chatRoom.Clients = append(chatRoom.Clients, client)
	client.ChatRoom = chatRoom
	chatRoom.Lock.Unlock()

	h.broadcast(chatRoom, client.Name, "has joined the chat room.") // it sends message to the server instead to chat

	fmt.Fprintf(client.Conn, "You joined chat room '%s'.\n", name)
}

func (h *ChatMessageHandler) leaveChatRoom(client *core.Client) {
	if client.ChatRoom == nil {
		return
	}

	client.ChatRoom.Lock.Lock()
	defer client.ChatRoom.Lock.Unlock()

	clients := client.ChatRoom.Clients
	for i, c := range clients {
		if c == client {
			client.ChatRoom.Clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	h.broadcast(client.ChatRoom, client.Name, "has joined the chat room.") // it sends message to the server instead to chat
	client.ChatRoom = nil
}

func (h *ChatMessageHandler) listUsers(client *core.Client) {
	if client.ChatRoom == nil {
		fmt.Fprintf(client.Conn, "You must be in a chat room to see the list of users.\n")
		return
	}

	client.ChatRoom.Lock.Lock()
	defer client.ChatRoom.Lock.Unlock()

	fmt.Fprintf(client.Conn, "Users in '%s':\n", client.ChatRoom.Name)
	for _, user := range client.ChatRoom.Clients {
		fmt.Fprintf(client.Conn, "- %s\n", user.Name)
	}
}
