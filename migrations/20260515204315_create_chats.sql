-- +goose Up
CREATE TABLE chats (
    id         VARCHAR(20) PRIMARY KEY,
    chat_id    INTEGER NOT NULL UNIQUE,
    title 	   TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE chats;

