package server

import (
	"log"
	"net/http"

	"github.com/devkekops/ctf-kartochki2-backend/internal/app/config"
	"github.com/devkekops/ctf-kartochki2-backend/internal/app/handlers"
	"github.com/devkekops/ctf-kartochki2-backend/internal/app/storage"
)

func Serve(cfg *config.Config) error {
	wordRepo, err := storage.NewWordRepo(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	baseHandler := handlers.NewBaseHandler(wordRepo, cfg.SecretKey, cfg.LicenseKeyHash)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: baseHandler,
	}

	return server.ListenAndServe()
}
