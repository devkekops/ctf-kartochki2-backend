package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/devkekops/ctf-kartochki2-backend/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type BaseHandler struct {
	mux       *chi.Mux
	fs        http.Handler
	secretKey string
	wordRepo  storage.WordRepository
}

func NewBaseHandler(wordRepo storage.WordRepository, secretKey string) *chi.Mux {
	cwd, _ := os.Getwd()
	root := filepath.Join(cwd)
	fs := http.FileServer(http.Dir(root))

	bh := &BaseHandler{
		mux:       chi.NewMux(),
		fs:        fs,
		secretKey: secretKey,
		wordRepo:  wordRepo,
	}

	bh.mux.Use(middleware.Logger)
	bh.mux.Use(authHandle(bh.secretKey))

	bh.mux.Handle("/*", fs)
	bh.mux.Route("/api", func(r chi.Router) {
		r.Get("/getFreeWords", bh.getFreeWords())
		r.Get("/getPaidWords", bh.getPaidWords())
	})

	return bh.mux
}

func (bh *BaseHandler) getFreeWords() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		words, err := bh.wordRepo.GetWords()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		buf, err := json.Marshal(words)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(buf)
		if err != nil {
			log.Println(err)
		}

	}
}

func (bh *BaseHandler) getPaidWords() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		words, err := bh.wordRepo.GetWords()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		buf, err := json.Marshal(words)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(buf)
		if err != nil {
			log.Println(err)
		}

	}
}
