package folder

import (
	"net/http"
	errs "svitlopus/internal/errors"

	"github.com/go-chi/chi/v5"
)

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) error {
	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(ErrEmptyFolderID, "Folder id is required")
	}

	ctx := r.Context()
	if err := h.folderService.DeleteFolder(ctx, folderId); err != nil {
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return nil
}
