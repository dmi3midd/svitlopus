package file

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoFile = errors.New("file not found in the repository")
)

// FileRepository defines the interface for file repository operations.
type FileRepository interface {
	// GetById returns a file by its ID.
	// It returns the ErrNoFile error if the file does not exist.
	GetById(ctx context.Context, id string) (*File, error)
	// GetByParentId returns a list of files by their parent ID.
	// It returns an empty slice if no files are found with the given parent ID.
	GetByParentId(ctx context.Context, parentId string, limit int, offset int) ([]File, error)
	// GetByTitleAndParentId returns a file by its title and parent ID.
	// It returns the ErrNoFile error if the file does not exist.
	GetByTitleAndParentId(ctx context.Context, title string, parentId string) (*File, error)
	// Create adds a new file to the database.
	Create(ctx context.Context, file *File) (*File, error)
	// Update modifies an existing file in the database.
	// It returns the ErrNoFile error if the file does not exist.
	Update(ctx context.Context, file *File) (*File, error)
	// Delete removes a file from the database by its ID.
	Delete(ctx context.Context, id string) error
}

type fileRepository struct {
	db *sqlx.DB
}

func NewFileRepo(db *sqlx.DB) FileRepository {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) GetById(ctx context.Context, id string) (*File, error) {
	op := "FileRepository.GetById"
	query := `
	SELECT id, title, mime, size, parent_id, file_id, message_id, created_at, updated_at 
	FROM files 
	WHERE id = $1
	`
	var file File
	err := r.db.GetContext(ctx, &file, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoFile)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &file, nil
}

func (r *fileRepository) GetByParentId(ctx context.Context, parentId string, limit int, offset int) ([]File, error) {
	op := "FileRepository.GetByParentId"
	query := `
	SELECT id, title, mime, size, parent_id, file_id, message_id, created_at, updated_at 
	FROM files 
	WHERE parent_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	var files []File
	err := r.db.SelectContext(ctx, &files, query, parentId, limit, offset)
	if err != nil {
		return []File{}, fmt.Errorf("%s: %w", op, err)
	}
	return files, nil
}

func (r *fileRepository) GetByTitleAndParentId(ctx context.Context, title string, parentId string) (*File, error) {
	op := "FileRepository.GetByTitleAndParentId"
	query := `
	SELECT id, title, mime, size, parent_id, file_id, message_id, created_at, updated_at 
	FROM files 
	WHERE title = $1 AND parent_id = $2
	`
	var file File
	err := r.db.GetContext(ctx, &file, query, title, parentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoFile)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &file, nil
}

func (r *fileRepository) Create(ctx context.Context, file *File) (*File, error) {
	op := "FileRepository.Create"
	query := `
	INSERT INTO files (id, title, mime, size, parent_id, file_id, message_id, created_at, updated_at) 
	VALUES (:id, :title, :mime, :size, :parent_id, :file_id, :message_id, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, file)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return file, nil
}

func (r *fileRepository) Update(ctx context.Context, file *File) (*File, error) {
	op := "FileRepository.Update"
	query := `
	UPDATE files 
	SET title = :title, parent_id = :parent_id, updated_at = :updated_id
	WHERE id = :id
	`
	result, err := r.db.NamedExecContext(ctx, query, file)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNoFile)
	}
	return file, nil
}

func (r *fileRepository) Delete(ctx context.Context, id string) error {
	op := "FileRepository.Delete"
	query := `
	DELETE FROM files 
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
