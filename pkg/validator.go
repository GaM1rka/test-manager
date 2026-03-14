package validator

import "errors"

func ValidateTodo(title string) error {
	if title == "" {
		return errors.New("Title cannot be empty")
	}
	return nil
}
