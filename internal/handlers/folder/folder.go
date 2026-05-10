package folder

import (
	"errors"
	"svitlopus/internal/models"
	"svitlopus/internal/services"
)

var (
	ErrEmptyFolderID  = errors.New("folder id is required")
	ErrFolderNotFound = errors.New("folder not found")
)

type FolderHandler struct {
	folderService services.FolderService
}

func NewFolderHandler(folderService services.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}

type FolderResponse struct {
	Folder models.Folder `json:"folder"`
}

type GetSubfoldersResponse struct {
	Folders []models.Folder `json:"folders"`
}

type CreateFolderRequest struct {
	Title    string `json:"title"`
	ParentId string `json:"parentId"`
}

type RenameFolderRequest struct {
	NewTitle string `json:"newTitle"`
}

type MoveFolderRequest struct {
	NewParentId string `json:"newParentId"`
}
