package auth

import (
	"errors"
	"fmt"
)

type FieldRequiredError struct {
	Field string
}

var (
	ErrIncorrectCredentials = errors.New("Incorrect credentials")
	ErrInvalidDataType      = errors.New("Invalid data type")
	ErrInternalServer       = errors.New("Internal server error")
	ErrUserNotInserted      = errors.New("User not inserted")
)

func (e *FieldRequiredError) Error() string {
	return fmt.Sprintf("This field is required: %s", e.Field)
}

func ErrFieldRequired(field string) error {
	return &FieldRequiredError{Field: field}
}
