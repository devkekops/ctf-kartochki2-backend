package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/devkekops/ctf-kartochki2-backend/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type LicenseKey struct {
	LicenseKey string `json:"licenseKey"`
}

type BaseHandler struct {
	mux            *chi.Mux
	fs             http.Handler
	secretKey      string
	licenseKeyHash string
	wordRepo       storage.WordRepository
}

func NewBaseHandler(wordRepo storage.WordRepository, secretKey string, licenseKeyHash string) *chi.Mux {
	cwd, _ := os.Getwd()
	root := filepath.Join(cwd)
	fs := http.FileServer(http.Dir(root))

	bh := &BaseHandler{
		mux:            chi.NewMux(),
		fs:             fs,
		secretKey:      secretKey,
		licenseKeyHash: licenseKeyHash,
		wordRepo:       wordRepo,
	}

	bh.mux.Use(middleware.Logger)
	bh.mux.Use(authHandle(bh.secretKey))

	bh.mux.Handle("/*", fs)
	bh.mux.Route("/api", func(r chi.Router) {
		r.Get("/getFreeWords", bh.getFreeWords())
		r.Get("/getPaidWords", bh.getPaidWords())
		r.Post("/activate", bh.activate())
	})

	return bh.mux
}

func (bh *BaseHandler) getFreeWords() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		words, err := bh.wordRepo.GetFreeWords()
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
		sessionCtx := req.Context().Value(sessionKey)
		session := sessionCtx.(Session)

		if session.Tariff != "paid" {
			http.Error(w, "You didn’t pay for the subscription!", http.StatusForbidden)
			return
		}

		words, err := bh.wordRepo.GetPaidWords()
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

func (bh *BaseHandler) activate() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionCtx := req.Context().Value(sessionKey)
		session := sessionCtx.(Session)

		var licenseKey LicenseKey
		if err := json.NewDecoder(req.Body).Decode(&licenseKey); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		hashBytes := sha256.Sum256([]byte(licenseKey.LicenseKey))
		hash := hex.EncodeToString(hashBytes[:])

		if hash != bh.licenseKeyHash {
			http.Error(w, "Wrong license key!", http.StatusUnprocessableEntity)
			return
		}

		session.Tariff = paidTariff
		session = createSession(bh.secretKey, session.UserID, session.Tariff)
		buf, err := json.Marshal(session)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		var sessionEnc = base64.StdEncoding.EncodeToString([]byte(buf))

		cookie := &http.Cookie{
			Name:  cookieName,
			Value: sessionEnc,
			Path:  cookiePath,
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
	}
}
