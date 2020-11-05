package handlers

import "net/http"

func (s *Server) registerUser() http.HandlerFunc	{
	return func(writer http.ResponseWriter, request *http.Request) {
		user, err := s.Domain.Register()
	}
}