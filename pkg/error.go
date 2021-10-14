package domain

import (
	"errors"
	"fmt"
)

// Error represents an error within the context of the service.
type Error struct {
	// Code is a machine-readable code.
	Code string `json:"code"`
	// Message is a human-readable message.
	Message string `json:"message"`
	// Inner is a wrapped error that is never shown to API consumers.
	Inner error `json:"-"`
}

// Error returns the string representation of the error message.
func (e Error) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s %s: %v", e.Code, e.Message, e.Inner)
	}
	return fmt.Sprintf("%s %s", e.Code, e.Message)
}

// Unwrap returns an inner error if any.
// It allows to use errors.Is() with eth.Error type.
func (e Error) Unwrap() error {
	return e.Inner
}

// ErrorCode returns the code of the error, if available.
func ErrorCode(err error) string {
	var e Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}
