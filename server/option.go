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
