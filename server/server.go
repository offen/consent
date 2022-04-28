package consent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
)

// NewHandler returns a http.Handler that serves the consent server using
// the given options.
func NewHandler(options ...Option) (http.Handler, error) {
	s := newDefaultServer()
	for _, option := range options {
		option(s)
	}
	return s, nil
}

func newDefaultServer() *server {
	return &server{
		logger:            log.New(io.Discard, "", log.Ldate),
		userCookieName:    "user",
		consentCookieName: "consent",
		userIDFunc: func() (string, error) {
			identifier := make([]byte, 16)
			if _, err := rand.Read(identifier); err != nil {
				return "", fmt.Errorf("userIDFunc: error reading random bytes. %w", err)
			}
			return hex.EncodeToString(identifier), nil
		},
	}
}

type server struct {
	logger            *log.Logger
	userCookieName    string
	consentCookieName string
	userIDFunc        func() (string, error)
}

// ServeHTTP handles a HTTP request
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := s.userIDFunc()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  s.userCookieName,
			Value: id,
		})
		w.Write([]byte("OK"))
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
