package server

import (
	"crypto/tls"
	"time"

	"github.com/emersion/go-smtp"
)

type Config struct {
	Addr              string
	Domain            string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxMessageBytes   int
	MaxRecipients     int
	CertFile          string
	KeyFile           string
	Endpoint          string
	ApiKey            string
	Username          string
	Password          string
	AllowInsecureAuth bool
}

func (c *Config) CreateServer() (*smtp.Server, error) {
	b := &Backend{
		Endpoint: c.Endpoint,
		ApiKey:   c.ApiKey,
		Username: c.Username,
		Password: c.Password,
	}

	s := smtp.NewServer(b)

	var tlsConfig *tls.Config

	s.AllowInsecureAuth = true

	if c.CertFile != "" && c.KeyFile != "" {
		tlsConfig = &tls.Config{}
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		if err != nil {
			return nil, err
		}
		s.TLSConfig = tlsConfig
		s.AllowInsecureAuth = c.AllowInsecureAuth
	}

	s.Addr = c.Addr
	s.Domain = c.Domain
	s.ReadTimeout = c.ReadTimeout
	s.WriteTimeout = c.WriteTimeout
	s.MaxMessageBytes = c.MaxMessageBytes
	s.MaxRecipients = c.MaxRecipients

	return s, nil
}
