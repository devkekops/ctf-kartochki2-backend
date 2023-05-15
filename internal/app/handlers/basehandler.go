package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devkekops/ctf-kartochki2-backend/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type BaseHandler struct {
	mux      *chi.Mux
	wordRepo storage.WordRepository
}

func NewBaseHandler(wordRepo storage.WordRepository) *chi.Mux {
	bh := &BaseHandler{
		mux:      chi.NewMux(),
		wordRepo: wordRepo,
	}

	bh.mux.Route("/api", func(r chi.Router) {
		r.Get("/getWords", bh.getWords())
	})

	return bh.mux
}

func (bh *BaseHandler) getWords() http.HandlerFunc {
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
