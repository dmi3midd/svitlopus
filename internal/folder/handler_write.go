package folder

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	errs "svitlopus/internal/errors"

	"github.com/go-chi/chi/v5"
)

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) error {
	var req CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewInternalServerError(err)
	}

	ctx := r.Context()
	folder, err := h.folderService.CreateFolder(ctx, req.Title, req.ParentId)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Parent folder not found",
			)
		}
		if errors.Is(err, ErrFolderAlreadyExist) {
			return errs.NewConflictError(
				ErrFolderAlreadyExist,
				"Folder with the same title already exists in the current directory",
			)
		}
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := FolderResponse{
		Folder: *folder,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}

func (h *FolderHandler) RenameFolder(w http.ResponseWriter, r *http.Request) error {
	var req RenameFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewInternalServerError(err)
	}

	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(fmt.Errorf("folder id is required"), "Folder id is required")
	}

	ctx := r.Context()
	folder, err := h.folderService.RenameFolder(ctx, folderId, req.NewTitle)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Folder not found",
			)
		}
		if errors.Is(err, ErrFolderAlreadyExist) {
			return errs.NewConflictError(
				ErrFolderAlreadyExist,
				"Folder with the same title already exists in the current directory",
			)
		}
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := FolderResponse{
		Folder: *folder,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}

func (h *FolderHandler) MoveFolder(w http.ResponseWriter, r *http.Request) error {
	var req MoveFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewInternalServerError(err)
	}

	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(fmt.Errorf("folder id is required"), "Folder id is required")
	}

	ctx := r.Context()
	folder, err := h.folderService.MoveFolder(ctx, folderId, req.NewParentId)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Folder or new parent folder not found",
			)
		}
		if errors.Is(err, ErrFolderAlreadyExist) {
			return errs.NewConflictError(
				ErrFolderAlreadyExist,
				"Folder with the same title already exists in the new parent directory",
			)
		}
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := FolderResponse{
		Folder: *folder,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) error {
	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(fmt.Errorf("folder id is required"), "Folder id is required")
	}

	ctx := r.Context()
	if err := h.folderService.DeleteFolder(ctx, folderId); err != nil {
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return nil
}
