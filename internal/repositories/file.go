package repositories

import (
	"context"
	"svitlopus/internal/models"
)

// FileRepository defines the interface for file repository operations.
type FileRepository interface {
	// GetById returns a file by its ID.
	// It returns the ErrFileNotFound error if the file does not exist.
	GetById(ctx context.Context, id string) (*models.File, error)
	// GetByParentId returns a list of files by their parent ID.
	// It returns an empty slice if no files are found with the given parent ID.
	GetByParentId(ctx context.Context, parentId string) ([]models.File, error)
	// Create adds a new file to the database.
	Create(ctx context.Context, file *models.File) error
	// Update modifies an existing file in the database.
	// It returns the ErrFileNotFound error if the file does not exist.
	Update(ctx context.Context, file *models.File) error
	// Delete removes a file from the database by its ID.
	Delete(ctx context.Context, id string) error
}
