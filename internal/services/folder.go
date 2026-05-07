package services

import (
	"context"
	"errors"
	"svitlopus/internal/models"
)

var (
	ErrFolderAlreadyExist = errors.New("folder already exist in current directory")
	ErrFolderNotFound     = errors.New("folder not found")
)

type FolderService interface {
	// GetFolder returns folder information.
	// Returns ErrFolderNotFound if no folder is found.
	GetFolder(ctx context.Context, folderId string) (*models.Folder, error)
	// GetSubfolders returns subfolders with pagination.
	// Returns ErrFolderNotFound if no folder is found.
	GetSubfolders(ctx context.Context, folderId string, limit int, offset int) ([]models.Folder, error)
	// GetFiles returns files in folder with pagination.
	// Returns ErrFolderNotFound if no folder is found.
	GetFiles(ctx context.Context, folderId string, limit int, offset int) ([]models.File, error)
	// CreateFolder creates and returns a new folder.
	// Returns ErrFolderAlreadyExist if a folder with the same title already exists in the current directory.
	CreateFolder(ctx context.Context, title, parentId string) (*models.Folder, error)
	// RenameFolder renames folder and returns the modified folder.
	// Returns ErrFolderNotFound if no folder is found.
	// Returns ErrFolderAlreadyExist if a folder with the same title already exists in the current directory.
	RenameFolder(ctx context.Context, id, newTitle string) (*models.Folder, error)
	// MoveFolder moves folder by changing it's parentId.
	// Returns ErrFolderNotFound if no folder is found.
	MoveFolder(ctx context.Context, id, parentId string) (*models.Folder, error)
	// DeleteFolder removes folder.
	DeleteFolder(ctx context.Context, id string)
}
