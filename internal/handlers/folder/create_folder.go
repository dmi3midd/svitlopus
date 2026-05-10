package folder

import (
	"encoding/json"
	"errors"
	"net/http"
	errs "svitlopus/internal/errors"
	"svitlopus/internal/services"
)

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) error {
	var req CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewInternalServerError(err)
	}

	ctx := r.Context()
	folder, err := h.folderService.CreateFolder(ctx, req.Title, req.ParentId)
	if err != nil {
		if errors.Is(err, services.ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Parent folder not found",
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
