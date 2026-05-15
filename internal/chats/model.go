package chats

import "time"

type Chat struct {
	Id        string    `json:"id" db:"id"`
	ChatId    int       `json:"chatId" db:"chat_id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
