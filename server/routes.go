package server

import "net/http"

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/ping", s.HandlePong()).Methods(http.MethodGet)

	// User routes
	s.router.HandleFunc("/users", s.HandleCreateUser()).Methods(http.MethodPost, http.MethodOptions)
	s.router.HandleFunc("/users", s.HandleGetUser()).Methods(http.MethodGet, http.MethodOptions)
	s.router.HandleFunc("/login", s.HandleLogin()).Methods(http.MethodPost, http.MethodOptions)
}

func (s *Server) HandlePong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(
			w,
			"pong",
			http.StatusOK,
			nil,
		)
	}
}
