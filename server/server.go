package consent

import "net/http"

func NewHandler(options ...Option) (http.Handler, error) {
	c := &config{}
	for _, option := range options {
		option(c)
	}
	return &server{
		config: c,
	}, nil
}

type server struct {
	config *config
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
