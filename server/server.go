package consent

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	case http.MethodPost:
		id, err := s.userIDFunc()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(
			w,
			s.makeCookie(s.userCookieName, id, time.Now().Add(s.cookieTTL)),
		)
		res := &response{
			Ok: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	case http.MethodGet:
		c, err := r.Cookie(s.consentCookieName)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("unable to find a cookie named '%s'", s.consentCookieName),
				http.StatusBadRequest,
			)
			return
		}

		raw, err := url.QueryUnescape(c.Value)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("unable to decode value of cookie named '%s'", s.consentCookieName),
				http.StatusBadRequest,
			)
			return
		}

		values, err := url.ParseQuery(raw)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("unable to deserialize value of cookie named '%s'", s.consentCookieName),
				http.StatusBadRequest,
			)
			return
		}

		var normalizedValues = make(map[string]string)
		for key, value := range values {
			normalizedValues[key] = value[0]
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"decisions": normalizedValues,
		})
	case http.MethodDelete:
		http.SetCookie(
			w,
			s.makeCookie(s.consentCookieName, "", time.Now().Add(-time.Hour)),
		)
		res := &response{
			Ok: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (s *server) makeCookie(name, value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		Path:     s.cookiePath,
		Domain:   s.cookieDomain,
		HttpOnly: true,
		Secure:   s.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

const (
	defaultUserCookieName    = "user"
	defaultConsentCookieName = "consent"
	defaultCookieTTL         = time.Hour * 24 * 31 * 6
	defaultCookieSecure      = true
)

func newDefaultServer() *server {
	return &server{
		logger:            log.New(io.Discard, "", log.Ldate),
		userCookieName:    defaultUserCookieName,
		consentCookieName: defaultConsentCookieName,
		cookieSecure:      defaultCookieSecure,
		cookieTTL:         defaultCookieTTL,
		userIDFunc: func() (string, error) {
			identifier := make([]byte, 16)
			if _, err := rand.Read(identifier); err != nil {
				return "", fmt.Errorf("userIDFunc: error reading random bytes. %w", err)
			}
			return hex.EncodeToString(identifier), nil
		},
	}
}
