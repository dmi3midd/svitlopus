package api

import (
	"net/http"
	"svitlopus/internal/config"
	"svitlopus/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	cfg *config.Http
	db  database.DBService
}

func NewServer(cfg *config.Http, db database.DBService) *http.Server {
	s := &Server{
		cfg: cfg,
		db:  db,
	}

	router := s.RegisterRoutes()
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
