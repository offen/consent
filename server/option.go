package consent

import "log"

// Option is a function used to configured the server.
type Option func(*server)

// WithLogger overrides the server's default logger with the given
// implementation.
func WithLogger(l *log.Logger) Option {
	return func(s *server) {
		s.logger = l
	}
}

// WithUserCookieName sets the name of the cookie that is used for storing
// user's consent decisions
func WithUserCookieName(n string) Option {
	return func(s *server) {
		s.userCookieName = n
	}
}

// WithConsentCookieName sets the name of the cookie that is used for storing
// user's consent decisions
func WithConsentCookieName(n string) Option {
	return func(s *server) {
		s.consentCookieName = n
	}
}

// WithUserIDFunc sets the function used for generating unique user identifiers
func WithUserIDFunc(f func() (string, error)) Option {
	return func(s *server) {
		s.userIDFunc = f
	}
}
