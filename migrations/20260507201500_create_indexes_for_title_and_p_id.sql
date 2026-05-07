-- +goose Up
CREATE UNIQUE INDEX idx_folders_title_parent ON folders(title, parent_id);
CREATE UNIQUE INDEX idx_files_title_parent ON files(title, parent_id);

-- +goose Down
DROP INDEX idx_folders_title_parent;
DROP INDEX idx_files_title_parent;
