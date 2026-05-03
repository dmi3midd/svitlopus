package api

import (
	"net/http"
	"svitlopus/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	cfg *config.Http
}

func NewServer(cfg *config.Http) *http.Server {
	s := &Server{
		cfg: cfg,
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
