// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"log"
	"time"
)

// Option is a function used to configured the server.
type Option func(*server)

// WithLogger overrides the server's default logger with the given
// implementation.
func WithLogger(l *log.Logger) Option {
	return func(s *server) {
		s.logger = l
	}
}

// WithCookieName sets the name of the cookie that is used for storing
// user's consent decisions
func WithCookieName(n string) Option {
	return func(s *server) {
		s.cookieName = n
	}
}

// WithCookiePath sets the Path attribute used when setting cookie headers.
func WithCookiePath(p string) Option {
	return func(s *server) {
		s.cookiePath = p
	}
}

// WithCookieDomain sets the Domain attribute used when setting cookie headers.
func WithCookieDomain(d string) Option {
	return func(s *server) {
		s.cookieDomain = d
	}
}

// WithCookieTTL defines the expected lifetime of a cookie.
func WithCookieTTL(d time.Duration) Option {
	return func(s *server) {
		s.cookieTTL = d
	}
}

// WithCookieSecure defines whether used cookies are using the Secure attribute
func WithCookieSecure(a bool) Option {
	return func(s *server) {
		s.cookieSecure = a
	}
}
