// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ServeHTTP handles a HTTP request
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthz":
		w.Write([]byte("OK"))
	case "/client.js":
		h.handleClientScript(w, r)
	case "/proxy":
		h.handleProxyHost(w, r)
	case "/consent":
		h.handleConsentRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *handler) handleProxyHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	if err := h.tpl.Execute(w, h.templateData); err != nil {
		http.Error(
			w,
			fmt.Sprintf("error rendering template: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *handler) handleClientScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/javascript")
	w.Write(h.clientScript)
}

type payload struct {
	Decisions map[scope]interface{} `json:"decisions"`
}

func (h *handler) handleConsentRequest(w http.ResponseWriter, r *http.Request) {
	d := decisions{}
	referrerURL, err := url.Parse(r.Referer())
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error parsing Referer header: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	referrer := domain(referrerURL.Host)

	if c, _ := r.Cookie(h.cookieName); c != nil {
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
		if err := json.NewEncoder(w).Encode(payload{Decisions: d[referrer]}); err != nil {
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
		d.update(&decisions{referrer: body.Decisions})

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
			h.makeCookie(h.cookieName, url.QueryEscape(encodedBody), time.Now().Add(h.cookieTTL)),
		)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&payload{Decisions: d[referrer]}); err != nil {
			http.Error(
				w,
				fmt.Sprintf("error encoding response payload: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
	case http.MethodDelete:
		http.SetCookie(
			w,
			h.makeCookie(h.cookieName, "", time.Now().Add(-time.Hour)),
		)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
