package folder

import (
	"encoding/json"
	"errors"
	"net/http"
	errs "svitlopus/internal/errors"
	"svitlopus/internal/services"

	"github.com/go-chi/chi/v5"
)

func (h *FolderHandler) GetFolder(w http.ResponseWriter, r *http.Request) error {
	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(ErrEmptyFolderID, "Folder id is required")
	}

	ctx := r.Context()
	folder, err := h.folderService.GetFolder(ctx, folderId)
	if err != nil {
		if errors.Is(err, services.ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Folder not found",
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
