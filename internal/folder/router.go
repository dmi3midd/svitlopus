package folder

import (
	errs "svitlopus/internal/errors"

	"github.com/go-chi/chi/v5"
)

func NewFolderRouter(handler *FolderHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/folders/{id}", errs.ErrorHandler(handler.GetFolder))
	r.Get("/folders/{id}/subfolders", errs.ErrorHandler(handler.GetSubfolders))
	r.Post("/folders", errs.ErrorHandler(handler.CreateFolder))
	r.Put("/folders/{id}/rename", errs.ErrorHandler(handler.RenameFolder))
	r.Put("/folders/{id}/move", errs.ErrorHandler(handler.MoveFolder))
	r.Delete("/folders/{id}", errs.ErrorHandler(handler.DeleteFolder))

	return r
}
