package folder

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	ErrFolderAlreadyExist = errors.New("folder already exist in current directory")
	ErrFolderNotFound     = errors.New("folder not found")
	ErrInvalidPagination  = errors.New("limit must be > 0 , < 30 and offset must be >= 0")
)

type FolderService interface {
	// GetFolder returns folder information.
	// Returns ErrFolderNotFound if no folder is found.
	GetFolder(ctx context.Context, folderId string) (*Folder, error)
	// GetSubfolders returns subfolders with pagination.
	// Returns ErrFolderNotFound if no folder is found.
	GetSubfolders(ctx context.Context, folderId string, limit, offset int) ([]Folder, error)
	// CreateFolder creates and returns a new folder.
	// Returns ErrFolderNotFound if parent folder is not found.
	// Returns ErrFolderAlreadyExist if a folder with the same title already exists in the current directory.
	CreateFolder(ctx context.Context, title, parentId string) (*Folder, error)
	// RenameFolder renames folder and returns the modified folder.
	// Returns ErrFolderNotFound if no folder is not found.
	// Returns ErrFolderAlreadyExist if a folder with the same title already exists in the current directory.
	RenameFolder(ctx context.Context, id, newTitle string) (*Folder, error)
	// MoveFolder moves folder by changing it's parentId.
	// Returns ErrFolderNotFound if no folder is not found.
	// Returns ErrFolderAlreadyExist if a folder with the same title already exists in the current directory.
	MoveFolder(ctx context.Context, id, newParentId string) (*Folder, error)
	// DeleteFolder removes folder.
	DeleteFolder(ctx context.Context, id string) error
}

type folderService struct {
	folderRepo FolderRepository
}

func NewFolderService(folderRepo FolderRepository) FolderService {
	return &folderService{
		folderRepo: folderRepo,
	}
}

func (s *folderService) GetFolder(ctx context.Context, folderId string) (*Folder, error) {
	op := "FolderService.GetFolder"

	folder, err := s.folderRepo.GetById(ctx, folderId)
	if err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return folder, nil
}

func (s *folderService) GetSubfolders(ctx context.Context, folderId string, limit, offset int) ([]Folder, error) {
	op := "FolderService.GetSubfolders"

	if limit <= 0 || offset < 0 {
		return []Folder{}, fmt.Errorf("%s: %w", op, ErrInvalidPagination)
	}

	if limit > 31 {

	}

	if _, err := s.folderRepo.GetById(ctx, folderId); err != nil {
		if errors.Is(err, ErrNoFolder) {
			return []Folder{}, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return []Folder{}, fmt.Errorf("%s: %w", op, err)
	}

	folders, err := s.folderRepo.GetByParentId(ctx, folderId, limit, offset)
	if err != nil {
		return []Folder{}, fmt.Errorf("%s: %w", op, err)
	}
	return folders, nil
}

func (s *folderService) CreateFolder(ctx context.Context, title, parentId string) (*Folder, error) {
	op := "FolderService.CreateFolder"

	if _, err := s.folderRepo.GetById(ctx, parentId); err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := s.folderRepo.GetByTitleAndParentId(ctx, title, parentId); err == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrFolderAlreadyExist)
	} else if !errors.Is(err, ErrNoFolder) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	folder := Folder{
		Id:        xid.New().String(),
		Title:     title,
		ParentId:  parentId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := s.folderRepo.Create(ctx, &folder); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &folder, nil
}

func (s *folderService) DeleteFolder(ctx context.Context, id string) error {
	op := "FolderService.DeleteFolder"

	if err := s.folderRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *folderService) MoveFolder(ctx context.Context, id, newParentId string) (*Folder, error) {
	op := "FolderService.MoveFolder"

	if _, err := s.folderRepo.GetById(ctx, newParentId); err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	folder, err := s.folderRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if folder.ParentId == newParentId {
		return folder, nil
	}

	if _, err := s.folderRepo.GetByTitleAndParentId(ctx, folder.Title, newParentId); err == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrFolderAlreadyExist)
	} else if !errors.Is(err, ErrNoFolder) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	folder.ParentId = newParentId
	folder.UpdatedAt = time.Now()

	updatedFolder, err := s.folderRepo.Update(ctx, folder)
	if err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return updatedFolder, nil
}

func (s *folderService) RenameFolder(ctx context.Context, id, newTitle string) (*Folder, error) {
	op := "FolderService.RenameFolder"

	folder, err := s.folderRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if folder.Title == newTitle {
		return folder, nil
	}

	if _, err := s.folderRepo.GetByTitleAndParentId(ctx, newTitle, folder.ParentId); err == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrFolderAlreadyExist)
	} else if !errors.Is(err, ErrNoFolder) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	folder.Title = newTitle
	folder.UpdatedAt = time.Now()

	updatedFolder, err := s.folderRepo.Update(ctx, folder)
	if err != nil {
		if errors.Is(err, ErrNoFolder) {
			return nil, fmt.Errorf("%s: %w", op, ErrFolderNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return updatedFolder, nil
}
