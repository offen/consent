// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

type handler struct {
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

type templateData struct {
	Script          *template.JS
	Styles          *template.CSS
	Wording         wording
	CustomTemplates *map[string]template.HTML
}

type wording struct {
	Paragraph string
	Yes       string
	No        string
}

// NewHandler returns a http.Handler that serves the consent server configured
// to use the the given options.
func NewHandler(options ...Option) (http.Handler, error) {
	s, err := newDefaultHandler()
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

//go:embed proxy/index.go.html
var proxyHostTemplate string

//go:embed proxy/proxy.js
var proxyScript string

//go:embed client/client.js
var clientScript string

func newDefaultHandler() (*handler, error) {
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

	return &handler{
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
