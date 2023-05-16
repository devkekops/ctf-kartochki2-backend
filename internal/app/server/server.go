package server

import (
	"net/http"

	"github.com/devkekops/ctf-kartochki2-backend/internal/app/config"
	"github.com/devkekops/ctf-kartochki2-backend/internal/app/handlers"
	"github.com/devkekops/ctf-kartochki2-backend/internal/app/storage"
)

func Serve(cfg *config.Config) error {
	var wordRepo = storage.NewWordRepo()

	var baseHandler = handlers.NewBaseHandler(wordRepo, cfg.SecretKey)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: baseHandler,
	}

	return server.ListenAndServe()
}
