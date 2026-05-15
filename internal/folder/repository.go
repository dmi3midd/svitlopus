package folder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoFolder = errors.New("folder not found in the repository")
)

// FolderRepository defines the interface for folder repository operations.
type FolderRepository interface {
	// GetById returns a folder by its ID.
	// It returns the ErrNoFolder error if the folder does not exist.
	GetById(ctx context.Context, id string) (*Folder, error)
	// GetByParentId returns a list of folders by their parent ID.
	// It returns an empty slice if no folders are found with the given parent ID.
	GetByParentId(ctx context.Context, parentId string, limit int, offset int) ([]Folder, error)
	// GetByTitleAndParentId returns a folder by its title and parent ID.
	// It returns the ErrNoFolder error if the folder does not exist.
	GetByTitleAndParentId(ctx context.Context, title, parentId string) (*Folder, error)
	// Create adds a new folder to the database.
	Create(ctx context.Context, folder *Folder) (*Folder, error)
	// Update modifies an existing folder in the database.
	// It returns the ErrNoFolder error if the folder does not exist.
	Update(ctx context.Context, folder *Folder) (*Folder, error)
	// Delete removes a folder from the database by its ID.
	Delete(ctx context.Context, id string) error
}

type folderRepository struct {
	db *sqlx.DB
}

func NewFolderRepo(db *sqlx.DB) FolderRepository {
	return &folderRepository{
		db: db,
	}
}

func (r *folderRepository) GetById(ctx context.Context, id string) (*Folder, error) {
	op := "FolderRepository.GetById"
	query := `
	SELECT id, title, parent_id, created_at, updated_at 
	FROM folders 
	WHERE id = $1
	`
	var folder Folder
	err := r.db.GetContext(ctx, &folder, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoFolder)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &folder, nil
}

func (r *folderRepository) GetByParentId(ctx context.Context, parentId string, limit int, offset int) ([]Folder, error) {
	op := "FolderRepository.GetByParentId"
	query := `
	SELECT id, title, parent_id, created_at, updated_at 
	FROM folders 
	WHERE parent_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	var folders []Folder
	err := r.db.SelectContext(ctx, &folders, query, parentId, limit, offset)
	if err != nil {
		return []Folder{}, fmt.Errorf("%s: %w", op, err)
	}
	return folders, nil
}

func (r *folderRepository) GetByTitleAndParentId(ctx context.Context, title string, parentId string) (*Folder, error) {
	op := "FolderRepository.GetByTitleAndParentId"
	query := `
	SELECT id, title, parent_id, created_at, updated_at 
	FROM folders 
	WHERE title = $1 AND parent_id = $2
	`
	var folder Folder
	err := r.db.GetContext(ctx, &folder, query, title, parentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoFolder)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &folder, nil
}

func (r *folderRepository) Create(ctx context.Context, folder *Folder) (*Folder, error) {
	op := "FolderRepository.Create"
	query := `
	INSERT INTO folders (id, title, parent_id, created_at, updated_at) 
	VALUES (:id, :title, :parent_id, :created_at, :updated_at)
	`
	if _, err := r.db.NamedExecContext(ctx, query, folder); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return folder, nil
}

func (r *folderRepository) Update(ctx context.Context, folder *Folder) (*Folder, error) {
	op := "FolderRepository.Update"
	query := `
	UPDATE folders 
	SET title = :title, parent_id = :parent_id, updated_at = :updated_at
	WHERE id = :id
	`
	result, err := r.db.NamedExecContext(ctx, query, folder)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrNoFolder)
	}
	return folder, nil
}

func (r *folderRepository) Delete(ctx context.Context, id string) error {
	op := "FolderRepository.Delete"
	query := `
	DELETE FROM folders 
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
