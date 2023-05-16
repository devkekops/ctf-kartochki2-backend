package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type key string

var (
	defaultTariff     = "free"
	paidTariff        = "paid"
	cookieName        = "session"
	cookiePath        = "/"
	sessionKey    key = "sessionKey"
)

type Session struct {
	UserID    string `json:"userID"`
	Tariff    string `json:"tariff"` //free or paid
	Signature string `json:"signature"`
}

func createSession(secretKey string, userID string, tariff string) Session {
	if userID == "" {
		userID = uuid.New().String()
	}
	if tariff == "" {
		tariff = defaultTariff
	}

	payload := []byte(userID + tariff)
	key := sha256.Sum256([]byte(secretKey))

	h := hmac.New(sha256.New, key[:])
	h.Write(payload[:])
	signature := hex.EncodeToString(h.Sum(nil))

	return Session{userID, tariff, signature}
}

func checkSignature(sessionEnc string, secretKey string) (Session, error) {
	sessionBuf, err := base64.StdEncoding.DecodeString(sessionEnc)
	if err != nil {
		return Session{}, err
	}

	var session Session
	err = json.Unmarshal(sessionBuf, &session)
	if err != nil {
		return Session{}, err
	}

	payload := []byte(session.UserID + session.Tariff)
	signature, err := hex.DecodeString(session.Signature)
	if err != nil {
		return Session{}, err
	}
	key := sha256.Sum256([]byte(secretKey))

	h := hmac.New(sha256.New, key[:])
	h.Write(payload[:])
	sign := h.Sum(nil)

	if hmac.Equal(sign, signature) {
		return session, nil
	} else {
		return Session{}, fmt.Errorf("invalid signature")
	}
}

func authHandle(secretKey string) (ah func(http.Handler) http.Handler) {
	ah = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var session Session
			sessionCookie, err := r.Cookie(cookieName)

			if err != nil {
				if err == http.ErrNoCookie {
					session = createSession(secretKey, "", "")
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
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
					log.Println(err)
					return
				}
			} else {
				sessionEnc := sessionCookie.Value
				session, err = checkSignature(sessionEnc, secretKey)
				if err != nil {
					session = createSession(secretKey, "", "")
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
				}
			}

			ctx := context.WithValue(r.Context(), sessionKey, session)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return
}
