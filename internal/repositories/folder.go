package repositories

import (
	"context"
	"svitlopus/internal/models"
)

// FolderRepository defines the interface for folder repository operations.
type FolderRepository interface {
	// GetById returns a folder by its ID.
	// It returns the ErrFolderNotFound error if the folder does not exist.
	GetById(ctx context.Context, id string) (*models.Folder, error)
	// GetByParentId returns a list of folders by their parent ID.
	// It returns an empty slice if no folders are found with the given parent ID.
	GetByParentId(ctx context.Context, parrentId string) ([]models.Folder, error)
	// Create adds a new folder to the database.
	Create(ctx context.Context, folder *models.Folder) (*models.Folder, error)
	// Update modifies an existing folder in the database.
	// It returns the ErrFolderNotFound error if the folder does not exist.
	Update(ctx context.Context, folder *models.Folder) (*models.Folder, error)
	// Delete removes a folder from the database by its ID.
	Delete(ctx context.Context, id string) error
}
