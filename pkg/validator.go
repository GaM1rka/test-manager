package validator

import "errors"

var (
	ErrInvalidTitle = errors.New("Title cannot be empty")
)

func ValidateTodo(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}
	return nil
}
