package chats

import (
	"context"
	"errors"
	"time"
)

var (
	ErrChatNotFound = errors.New("chat not found")
)

// TODO: Redesign logic with delay.
type ChatManager interface {
	// AddChat adds chat id to chat manager.
	AddChat(ctx context.Context, chatId int)
	// NextChat set current index on the next chat id.
	// It returns next chat id and time to wait before uploading next file.
	NextChat() (int, time.Duration)
}

type chatManager struct {
	chatIds      []int
	currentIndex int
}

func NewChatManager(chatIds []int) ChatManager {
	return &chatManager{
		chatIds:      chatIds,
		currentIndex: 0,
	}
}

func (m *chatManager) AddChat(ctx context.Context, chatId int) {
	m.chatIds = append(m.chatIds, chatId)
}

func (m *chatManager) NextChat() (int, time.Duration) {
	count := len(m.chatIds)
	if count == 0 {
		return 0, 0
	}

	selectedChat := m.chatIds[m.currentIndex]

	m.currentIndex = (m.currentIndex + 1) % count

	delayMilliseconds := max(3000/count, 40)
	delay := time.Duration(delayMilliseconds) * time.Millisecond

	return selectedChat, delay
}
