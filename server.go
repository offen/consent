// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type server struct {
	logger       *log.Logger
	cookieName   string
	cookieDomain string
	cookiePath   string
	cookieTTL    time.Duration
	cookieSecure bool
	tpl          *template.Template
	templateData templateData
	clientScript []byte
}

// ServeHTTP handles a HTTP request
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthz":
		w.Write([]byte("OK"))
	case "/client.js":
		s.handleClientScript(w, r)
	case "/proxy":
		s.handleProxyHost(w, r)
	case "/consent":
		s.handleConsentRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

// NewHandler returns a http.Handler that serves the consent server using
// the given options.
func NewHandler(options ...Option) (http.Handler, error) {
	s, err := newDefaultServer()
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(s)
	}
	return s, nil
}

type payload struct {
	Decisions decisions `json:"decisions"`
}

func (s *server) handleProxyHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	if err := s.tpl.Execute(w, s.templateData); err != nil {
		http.Error(
			w,
			fmt.Sprintf("error rendering template: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}
}

func (s *server) handleClientScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	w.Write(s.clientScript)
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
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload{Decisions: d}); err != nil {
			http.Error(
				w,
				fmt.Sprintf("error encoding response payload: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
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

//go:embed proxy/index.go.html
var proxyHostTemplate string

//go:embed proxy/proxy.js
var proxyScript string

//go:embed client/client.js
var clientScript string

func newDefaultServer() (*server, error) {
	tpl := template.New("proxy")
	if _, err := tpl.Parse(string(proxyHostTemplate)); err != nil {
		return nil, fmt.Errorf("newDefaultServer: error parsing template: %w", err)
	}
	minifiedProxyScript, err := minifyJS(proxyScript)
	if err != nil {
		return nil, fmt.Errorf("newDefaultServer: error minifying proxy script: %w", err)
	}
	safeScript := template.JS(minifiedProxyScript)

	minifiedClientScript, err := minifyJS(clientScript)
	if err != nil {
		return nil, fmt.Errorf("newDefaultServer: error minifying client script: %w", err)
	}

	return &server{
		logger:       log.New(io.Discard, "", log.Ldate),
		cookieName:   defaultConsentCookieName,
		cookieSecure: defaultCookieSecure,
		cookieTTL:    defaultCookieTTL,
		clientScript: []byte(minifiedClientScript),
		tpl:          tpl,
		templateData: templateData{
			Script: &safeScript,
		},
	}, nil
}

type templateData struct {
	Script  *template.JS
	Styles  *template.CSS
	Wording wording
}

type wording struct {
	Paragraph string
	Yes       string
	No        string
}
