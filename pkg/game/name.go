package game

import (
	"errors"
	"regexp"
)

var (
	// ErrEmptyName is returned when the name is empty.
	ErrEmptyName = errors.New("name is empty")
	// ErrNameTooLong is returned when the name is too long.
	ErrNameTooLong = errors.New("name is too long")
	// ErrInvalidName is returned when the name contains invalid characters.
	ErrInvalidName = errors.New("name contains invalid characters")
)

// nameRegex is a regular expression that matches a name.
var nameRegex = regexp.MustCompile(`^[a-z]+$`)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	if len(name) > 100 {
		return ErrNameTooLong
	}

	if !nameRegex.MatchString(name) {
		return ErrInvalidName
	}

	return nil
}
