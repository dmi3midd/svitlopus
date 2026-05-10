package worker

import "time"

type Worker struct {
	Id        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	BotToken  string    `json:"botToken" db:"bot_token"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	ChatId    int       `json:"chatId" db:"chat_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
