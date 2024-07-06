package chat

import (
	"sync"
	"tcp-chat/internal/core"
)

var (
	ChatRooms     = make(map[string]*core.ChatRoom)
	ChatRoomsLock sync.Mutex
)

func InitializeChatRooms() {}
