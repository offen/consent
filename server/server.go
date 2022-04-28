package consent

import (
	"io"
	"log"
	"net/http"
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

func newDefaultServer() *server {
	return &server{
		logger: log.New(io.Discard, "", log.Ldate),
	}
}

type server struct {
	logger *log.Logger
}

// ServeHTTP handles a HTTP request
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
