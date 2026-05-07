-- +goose Up
INSERT INTO folders (id, title, parent_id) VALUES ('root', 'root', 'root');

-- +goose Down
DELETE FROM folders WHERE id = 'root';
