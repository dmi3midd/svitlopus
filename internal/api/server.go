package api

import (
	"net/http"
	"svitlopus/internal/config"
	"svitlopus/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	cfg *config.Config
	db  database.DBService
}

func NewServer(cfg *config.Config, db database.DBService) *http.Server {
	s := &Server{
		cfg: cfg,
		db:  db,
	}

	router := s.RegisterRoutes()
	return &http.Server{
		Addr:         cfg.Http.Address,
		Handler:      router,
		IdleTimeout:  cfg.Http.IdleTimeout,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
	}
}
