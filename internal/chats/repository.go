package chats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRepoChatNotFound = errors.New("repository: chat not found")
)

type ChatRepository interface {
	// GetById returns a chat by its Id.
	// It returns the ErrRepoChatNotFound error if the chat does not exist.
	GetById(ctx context.Context, id string) (*Chat, error)
	// GetByChatId returns a chat by its Id.
	// It returns the ErrRepoChatNotFound error if the chat does not exist.
	GetByChatId(ctx context.Context, chatId int) (*Chat, error)
	// GetAll returns a list of chats.
	// It returns an empty slice if no chats are found.
	GetAll(ctx context.Context) ([]Chat, error)
	// Create adds a new chat to the database.
	Create(ctx context.Context, chat *Chat) (*Chat, error)
	// Delete removes a chat from the database by its ID.
	Delete(ctx context.Context, id string) error
}

type chatRepository struct {
	db *sqlx.DB
}

func NewChatRepo(db *sqlx.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (r *chatRepository) GetById(ctx context.Context, id string) (*Chat, error) {
	op := "ChatRepository.GetById"
	query := `
	SELECT id, chat_id, title, created_at 
	FROM chats 
	WHERE id = $1
	`
	var chat Chat
	if err := r.db.GetContext(ctx, &chat, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrRepoChatNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &chat, nil
}

func (r *chatRepository) GetByChatId(ctx context.Context, chatId int) (*Chat, error) {
	op := "ChatRepository.GetByChatId"
	query := `
	SELECT id, chat_id, title, created_at 
	FROM chats 
	WHERE chat_id = $1
	`
	var chat Chat
	if err := r.db.GetContext(ctx, &chat, query, chatId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrRepoChatNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &chat, nil
}

func (r *chatRepository) GetAll(ctx context.Context) ([]Chat, error) {
	op := "ChatRepository.GetAll"
	query := `
	SELECT id, chat_id, title, created_at 
	FROM chats
	`
	var chats []Chat
	if err := r.db.SelectContext(ctx, &chats, query); err != nil {
		return []Chat{}, fmt.Errorf("%s: %w", op, err)
	}
	return chats, nil
}

func (r *chatRepository) Create(ctx context.Context, chat *Chat) (*Chat, error) {
	op := "ChatRepository.Create"
	query := `
	INSERT INTO chats (id, chat_id, title, created_at)
	VALUES (:id, :chat_id, :title, :created_at)
	`
	if _, err := r.db.NamedExecContext(ctx, query, chat); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return chat, nil
}

func (r *chatRepository) Delete(ctx context.Context, id string) error {
	op := "ChatRepository.Delete"
	query := `
	DELETE FROM chats 
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
