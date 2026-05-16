package telegramadapter

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"gopkg.in/telebot.v4"
)

var (
	ErrUploadFailed = errors.New("failed to upload file to Telegram")
)

type TelegramClient interface {
	Ping() error
	UploadFile(file multipart.File, chatId int) (*UploadedFile, error)
}

type telegramClient struct {
	botToken string
	baseURL  string
	client   *http.Client
	bot      *telebot.Bot
}

func NewTelegramClient(botToken string, baseURL string) (TelegramClient, error) {
	pref := telebot.Settings{
		URL:    baseURL,
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Client: &http.Client{},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram Bot: %v", err)
	}
	return &telegramClient{
		botToken: botToken,
		baseURL:  baseURL,
		client:   &http.Client{},
		bot:      bot,
	}, nil
}

func (c *telegramClient) Ping() error {
	url := fmt.Sprintf("%s/bot%s/getMe", c.baseURL, c.botToken)
	resp, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to Telegram Bot API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram Bot API returned status code %d", resp.StatusCode)
	}

	return nil
}

type UploadedFile struct {
	TgMsgId  int
	TgFileId string
	Size     int
	MimeType string
}

func (c *telegramClient) UploadFile(file multipart.File, chatId int) (*UploadedFile, error) {
	op := "TelegramClient.UploadFile"
	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)
	mimeType := http.DetectContentType(buffer)

	chat := &telebot.Chat{
		ID: int64(chatId),
	}
	document := &telebot.Document{
		MIME: mimeType,
		File: telebot.FromReader(file),
	}
	msg, err := c.bot.Send(chat, document)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrUploadFailed)
	}
	return &UploadedFile{
		TgMsgId:  msg.ID,
		TgFileId: document.FileID,
		Size:     int(document.FileSize),
		MimeType: document.MIME,
	}, nil
}
