package consent

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

//go:embed proxy/proxy.js
var proxyScript []byte

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
	logger       *log.Logger
	cookieName   string
	cookieDomain string
	cookiePath   string
	cookieTTL    time.Duration
	cookieSecure bool
}

type payload struct {
	Decisions decisions `json:"decisions"`
}

func (s *server) handleProxyScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	w.Write(proxyScript)
}

func (s *server) handleConsentRequest(w http.ResponseWriter, r *http.Request) {
	d := decisions{}

	if c, _ := r.Cookie(s.cookieName); c != nil {
		raw, err := url.QueryUnescape(c.Value)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error unescaping cookie value: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}
		decisionsFromCookie, err := parseDecisions(raw)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error parsing unescaped cookie value: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}
		d.update(decisionsFromCookie)
	}

	switch r.Method {
	case http.MethodPost:
		body := payload{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(
				w,
				fmt.Sprintf("error decoding body: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}
		d.update(&body.Decisions)

		encodedBody, err := d.encode()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error encoding decisions as cookie: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		http.SetCookie(
			w,
			s.makeCookie(s.cookieName, url.QueryEscape(encodedBody), time.Now().Add(s.cookieTTL)),
		)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&payload{Decisions: d}); err != nil {
			http.Error(
				w,
				fmt.Sprintf("error encoding response payload: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload{Decisions: d}); err != nil {
			http.Error(
				w,
				fmt.Sprintf("error encoding response payload: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
	case http.MethodDelete:
		http.SetCookie(
			w,
			s.makeCookie(s.cookieName, "", time.Now().Add(-time.Hour)),
		)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

// ServeHTTP handles a HTTP request
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/proxy.js":
		s.handleProxyScript(w, r)
	case "/":
		s.handleConsentRequest(w, r)
	default:
		http.NotFound(w, r)
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
	defaultConsentCookieName = "consent"
	defaultCookieTTL         = time.Hour * 24 * 31 * 6
	defaultCookieSecure      = true
)

func newDefaultServer() *server {
	return &server{
		logger:       log.New(io.Discard, "", log.Ldate),
		cookieName:   defaultConsentCookieName,
		cookieSecure: defaultCookieSecure,
		cookieTTL:    defaultCookieTTL,
	}
}
