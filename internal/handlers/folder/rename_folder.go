package folder

import (
	"encoding/json"
	"errors"
	"net/http"
	errs "svitlopus/internal/errors"
	"svitlopus/internal/services"

	"github.com/go-chi/chi/v5"
)

func (h *FolderHandler) RenameFolder(w http.ResponseWriter, r *http.Request) error {
	var req RenameFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewInternalServerError(err)
	}

	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(ErrEmptyFolderID, "Folder id is required")
	}

	ctx := r.Context()
	folder, err := h.folderService.RenameFolder(ctx, folderId, req.NewTitle)
	if err != nil {
		if errors.Is(err, services.ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Folder not found",
			)
		}
		if errors.Is(err, services.ErrFolderAlreadyExist) {
			return errs.NewConflictError(
				services.ErrFolderAlreadyExist,
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
