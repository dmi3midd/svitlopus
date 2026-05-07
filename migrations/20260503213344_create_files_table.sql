-- +goose Up
CREATE TABLE files (
    id          VARCHAR(20) PRIMARY KEY,
    title       VARCHAR(48) NOT NULL,
    mime        TEXT NOT NULL,
    size        INTEGER NOT NULL,
    parent_id   VARCHAR(20) NOT NULL,
    file_id     TEXT NOT NULL,
    message_id  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES folders (id) ON DELETE CASCADE
);

CREATE INDEX idx_files_parent ON files(parent_id);

-- +goose Down
DROP TABLE files;
