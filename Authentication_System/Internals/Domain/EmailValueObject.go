package domain

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrEmptyEmail         = errors.New("email cannot be empty")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type Email struct {
	value string
}

func NewEmail(rawEmail string) (Email, error) {
	cleanEmail := strings.TrimSpace(strings.ToLower(rawEmail))
	
	if cleanEmail == "" {
		return Email{}, ErrEmptyEmail
	}

	if _, err := mail.ParseAddress(cleanEmail); err != nil {
		return Email{}, ErrInvalidEmailFormat
	}

	return Email{value: cleanEmail}, nil
}

func (e Email) Value() string {
	return e.value
}