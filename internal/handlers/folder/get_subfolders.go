package folder

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	errs "svitlopus/internal/errors"
	"svitlopus/internal/services"

	"github.com/go-chi/chi/v5"
)

func (h *FolderHandler) GetSubfolders(w http.ResponseWriter, r *http.Request) error {
	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		return errs.NewBadRequestError(ErrEmptyFolderID, "Folder id is required")
	}

	limit := 20
	offset := 0
	query := r.URL.Query()
	if l := query.Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if o := query.Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	ctx := r.Context()
	folders, err := h.folderService.GetSubfolders(ctx, folderId, limit, offset)
	if err != nil {
		if errors.Is(err, services.ErrFolderNotFound) {
			return errs.NewNotFoundError(
				ErrFolderNotFound,
				"Parent folder not found",
			)
		}
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := GetSubfoldersResponse{
		Folders: folders,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}
