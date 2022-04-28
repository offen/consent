package consent

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

type server struct {
	logger            *log.Logger
	userCookieName    string
	consentCookieName string
	cookieDomain      string
	cookiePath        string
	cookieTTL         time.Duration
	cookieSecure      bool
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
			Name:     s.userCookieName,
			Value:    id,
			Path:     s.cookiePath,
			Domain:   s.cookieDomain,
			Expires:  time.Now().Add(s.cookieTTL),
			HttpOnly: true,
			Secure:   s.cookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
		res := &response{
			Ok: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func newDefaultServer() *server {
	return &server{
		logger:            log.New(io.Discard, "", log.Ldate),
		userCookieName:    "user",
		consentCookieName: "consent",
		cookieSecure:      true,
		cookieTTL:         time.Hour * 24 * 31 * 6,
		userIDFunc: func() (string, error) {
			identifier := make([]byte, 16)
			if _, err := rand.Read(identifier); err != nil {
				return "", fmt.Errorf("userIDFunc: error reading random bytes. %w", err)
			}
			return hex.EncodeToString(identifier), nil
		},
	}
}
