package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/emersion/go-smtp"
)

// A Session is returned after EHLO.
type Session struct {
	backend *Backend
	from    string
	to      []string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)

	if err != nil {
		s.log("Failed to read data: ")

		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.backend.Endpoint, bytes.NewReader(b))

	if err != nil {
		s.log("Failed to create HTTP request: ")

		return err
	}

	if s.backend.ApiKey != "" {
		req.Header.Set("X-Api-Key", s.backend.ApiKey)
	}

	req.Header.Set("X-Mail-From", s.from)

	for _, to := range s.to {
		req.Header.Add("X-Rcpt-To", to)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		s.log("Failed to make HTTP request: ")

		return err
	}

	s.log(res.Status)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected response status: %d", res.StatusCode)
	}

	return nil
}

func (s *Session) log(message string) {
	log.Println(message, s.from, "=>", strings.Join(s.to, ", "))
}

func (s *Session) Reset() {
	s.from = ""
	s.to = []string{}
}

func (s *Session) Logout() error {
	s.backend = nil
	return nil
}
