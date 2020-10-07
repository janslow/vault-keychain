package vaultkeychain

import (
	"errors"
	"net/url"
)

var (
	ErrTokenNotFound = errors.New("no token found")
)

type Server struct {
	Address *url.URL
}

func (s *Server) Token() (string, error) {
	return read(s.Address)
}

func (s *Server) SetToken(token string) error {
	return write(s.Address, token)
}

func (s *Server) ClearToken() error {
	return clear(s.Address)
}
