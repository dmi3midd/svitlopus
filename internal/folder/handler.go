package folder

import (
	"errors"
)

var (
	ErrEmptyFolderID = errors.New("folder id is required")
	// ErrFolderNotFound = errors.New("folder not found")
)

type FolderHandler struct {
	folderService FolderService
}

func NewFolderHandler(folderService FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}

type FolderResponse struct {
	Folder Folder `json:"folder"`
}

type GetSubfoldersResponse struct {
	Folders []Folder `json:"folders"`
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
