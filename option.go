// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Option is a function used to configured the server.
type Option func(*server) error

// WithLogger overrides the server's default logger with the given
// implementation.
func WithLogger(l *log.Logger) Option {
	return func(s *server) error {
		s.logger = l
		return nil
	}
}

// WithCookieName sets the name of the cookie that is used for storing
// user's consent decisions
func WithCookieName(n string) Option {
	return func(s *server) error {
		s.cookieName = n
		return nil
	}
}

// WithCookiePath sets the Path attribute used when setting cookie headers.
func WithCookiePath(p string) Option {
	return func(s *server) error {
		s.cookiePath = p
		return nil
	}
}

// WithCookieDomain sets the Domain attribute used when setting cookie headers.
func WithCookieDomain(d string) Option {
	return func(s *server) error {
		s.cookieDomain = d
		return nil
	}
}

// WithCookieTTL defines the expected lifetime of a cookie.
func WithCookieTTL(d time.Duration) Option {
	return func(s *server) error {
		s.cookieTTL = d
		return nil
	}
}

// WithCookieSecure defines whether used cookies are using the Secure attribute
func WithCookieSecure(a bool) Option {
	return func(s *server) error {
		s.cookieSecure = a
		return nil
	}
}

// WithCustomizedWording passes custom copy to be used in the default consent UI
func WithCustomizedWording(copy, yes, no string) Option {
	return func(s *server) error {
		s.templateData.Wording.Paragraph = copy
		s.templateData.Wording.Yes = yes
		s.templateData.Wording.No = no
		return nil
	}
}

// WithStylesheet adds a stylesheet that is injected into the iframe element.
func WithStylesheet(loc string) Option {
	return func(s *server) error {
		if loc == "" {
			return nil
		}
		b, err := os.ReadFile(loc)
		if err != nil {
			return fmt.Errorf("WithCustomStyles: error reading given file: %w", err)
		}
		css := template.CSS(string(b))
		s.templateData.Styles = &css
		return nil
	}
}

// WithTemplatesDirectory configures the server to look for custom templates
// in the given location.
func WithTemplatesDirectory(dir string) Option {
	return func(s *server) error {
		if dir == "" {
			return nil
		}

		templates := map[string]template.HTML{}
		if err := filepath.WalkDir(dir, func(path string, di fs.DirEntry, err error) error {
			if di.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".html" {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("WithTemplatesDirectory: error reading file %s: %w", path, err)
			}
			id := filepath.Base(path)
			id = strings.TrimSuffix(id, ".html")
			templates[id] = template.HTML(string(b))
			return nil
		}); err != nil {
			return fmt.Errorf("WithTemplatesDirectory: error walking directory %s: %w", dir, err)
		}
		s.templateData.CustomTemplates = &templates
		return nil
	}
}
