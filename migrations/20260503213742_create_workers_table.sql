-- +goose Up
CREATE TABLE workers (
    id         VARCHAR(20) PRIMARY KEY,
    username   VARCHAR(48) NOT NULL UNIQUE,
    bot_token  TEXT NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    chat_id    INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE workers;
