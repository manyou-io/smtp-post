package server

import (
	"errors"

	"github.com/emersion/go-smtp"
)

type Backend struct {
	Endpoint string
	ApiKey   string
	Username string
	Password string
}

func (b *Backend) Login(_ *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username == b.Username && password == b.Password {
		return &Session{backend: b}, nil
	}

	return nil, errors.New("invalid username or password")
}

func (b *Backend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	if b.Username == "" && b.Password == "" {
		return &Session{backend: b}, nil
	}

	return nil, smtp.ErrAuthRequired
}
