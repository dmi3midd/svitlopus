package api

import (
	"encoding/json"
	"net/http"
	hfolder "svitlopus/internal/handlers/folder"
	"svitlopus/internal/repositories"
	"svitlopus/internal/routers"
	"svitlopus/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	folderRepo := repositories.NewFolderRepo(s.db.GetDB())
	folderService := services.NewFolderService(folderRepo)
	folderHandler := hfolder.NewFolderHandler(folderService)
	folderRouter := routers.NewFolderRouter(folderHandler)

	mainRouter := chi.NewRouter()

	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mainRouter.Use(middleware.RequestID)
	mainRouter.Use(middleware.Recoverer)

	mainRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.db.Health())
	})

	mainRouter.Group(func(r chi.Router) {
		r.Mount("/svitlopus", folderRouter)
	})

	return mainRouter
}
